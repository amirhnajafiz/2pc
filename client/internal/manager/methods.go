package manager

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/F24-CSE535/2pc/client/internal/utils"
	"github.com/F24-CSE535/2pc/client/pkg/models"
	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

func (m *Manager) Performance() string {
	return fmt.Sprintf("throughput: %f tps, latency: %f ms", utils.Average(m.throughput), utils.Average(m.latency))
}

func (m *Manager) PrintBalance(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// get the shard
	cluster, err := m.storage.GetClientShard(argv[0])
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
	}

	// get the cluster's services
	services := strings.Split(m.dialer.Nodes[fmt.Sprintf("E%s", cluster)], ":")

	// make RPC call
	output := fmt.Sprintf("--  server  -  %s  --\n", argv[0])
	for _, svc := range services {
		if balance, err := m.dialer.PrintBalance(svc, argv[0]); err != nil {
			return fmt.Errorf("server failed: %v", err).Error()
		} else {
			output = fmt.Sprintf("%s     %s     -  %d\n", output, svc, balance)
		}
	}

	return output
}

func (m *Manager) PrintLogs(argc int, argv []string) ([]string, string) {
	// check the number of arguments
	if argc < 1 {
		return nil, "not enough arguments"
	}

	// make RPC call
	if list, err := m.dialer.PrintLogs(argv[0]); err != nil {
		return nil, err.Error()
	} else {
		return list, ""
	}
}

func (m *Manager) PrintDatastore(argc int, argv []string) ([]string, string) {
	// check the number of arguments
	if argc < 1 {
		return nil, "not enough arguments"
	}

	// make RPC call
	list, err := m.dialer.PrintDatastore(argv[0])
	if err != nil {
		return nil, err.Error()
	}

	// create a list of sessions
	records := make([]string, 0)
	for _, msg := range list {
		if ses, err := m.storage.GetSessionById(int(msg.GetSessionId())); err == nil {
			records = append(records, fmt.Sprintf(
				"\t[<%d, %s>, (%s, %s, %d)]",
				msg.GetBallotNumberSequence(),
				msg.GetBallotNumberPid(),
				ses.Sender,
				ses.Receiver,
				ses.Amount,
			))
		} else {
			log.Printf("failed to get session %d: %v\n", msg.GetSessionId(), err)
		}
	}

	return records, ""
}

func (m *Manager) PrintDatastores(argc int, argv []string) ([]string, string) {
	// set a list of servers
	servers := strings.Split(m.dialer.Nodes["all"], ":")

	// create a list of records
	records := make([]string, 0)

	// loop over servers and get messages
	for _, svc := range servers {
		// make RPC call
		list, err := m.dialer.PrintDatastore(svc)
		if err != nil {
			return nil, err.Error()
		}

		// get sessions from the database
		output := fmt.Sprintf("%s:\n", svc)
		for _, msg := range list {
			if ses, err := m.storage.GetSessionById(int(msg.GetSessionId())); err == nil {
				output = fmt.Sprintf(
					"%s\t[<%d, %s>, (%s, %s, %d)]\n",
					output,
					msg.GetBallotNumberSequence(),
					msg.GetBallotNumberPid(),
					ses.Sender,
					ses.Receiver,
					ses.Amount,
				)
			} else {
				log.Printf("failed to get session %d: %v\n", msg.GetSessionId(), err)
			}
		}

		records = append(records, output)
	}

	return records, ""
}

func (m *Manager) Transaction(argc int, argv []string) (string, bool) {
	// check the number of arguments
	if argc < 3 {
		return "not enough arguments", false
	}

	// extract data from input command
	sender := argv[0]
	receiver := argv[1]
	amount, _ := strconv.Atoi(argv[2])
	sessionId := m.memory.GetSession()

	// get shards
	senderCluster, err := m.storage.GetClientShard(sender)
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error(), false
	}
	receiverCluster, err := m.storage.GetClientShard(receiver)
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error(), false
	}

	// create a new session
	session := models.Session{
		Id:       sessionId,
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
		Replys:   make([]*database.ReplyMsg, 0),
	}

	// check for inter or cross shard
	if senderCluster == receiverCluster {
		session.Type = "inter-shard"
		session.Participants = []string{senderCluster}

		// for inter-shard send request message to the cluster
		if err := m.dialer.Request(senderCluster, sender, receiver, amount, sessionId); err != nil {
			return fmt.Errorf("server failed: %v", err).Error(), false
		}
	} else {
		session.Type = "cross-shard"
		session.Participants = []string{senderCluster, receiverCluster}
		session.Acks = make([]*database.AckMsg, 0)

		// for cross-shard send prepare messages to both clusters
		if err := m.dialer.Prepare(senderCluster, sender, sender, receiver, amount, sessionId); err != nil {
			return fmt.Errorf("sender server failed: %v", err).Error(), false
		}
		if err := m.dialer.Prepare(receiverCluster, receiver, sender, receiver, amount, sessionId); err != nil {
			return fmt.Errorf("receiver server failed: %v", err).Error(), false
		}
	}

	// save the transaction into cache
	session.StartedAt = time.Now()
	m.cache[sessionId] = &session

	// store session for future optimizations
	if err := m.storage.InsertSession(&session); err != nil {
		log.Printf("failed to store cross-shard metric: %v\n", err)
	}

	return fmt.Sprintf("transaction %d (%s, %s, %d): sent", sessionId, sender, receiver, amount), true
}

func (m *Manager) RoundTrip(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// get the shard
	cluster, err := m.storage.GetClientShard(argv[0])
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
	}

	// make RPC call
	if err := m.dialer.Request(cluster, argv[0], "", 0, 0); err != nil {
		return fmt.Errorf("server failed: %v", err).Error()
	}

	return "roundtrip sent"
}

func (m *Manager) Block(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// make RPC call
	if err := m.dialer.Block(argv[0]); err != nil {
		return fmt.Errorf("server failed: %v", err).Error()
	}

	return "blocked"
}

func (m *Manager) Unblock(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// make RPC call
	if err := m.dialer.Unblock(argv[0]); err != nil {
		return fmt.Errorf("server failed: %v", err).Error()
	}

	return "unblocked"
}

func (m *Manager) LoadTests(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// load tests
	tc, err := utils.CSVParseTestcaseFile(argv[0])
	if err != nil {
		return err.Error()
	}

	// set tests
	m.tests = tc
	m.index = 0

	return fmt.Sprintf("load %d testsets", len(tc))
}

func (m *Manager) ShardsRebalance(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// open the file
	file, err := os.Open(argv[0])
	if err != nil {
		return err.Error()
	}
	defer file.Close()

	// create a scanner
	scanner := bufio.NewScanner(file)

	var (
		firstCluster string
		account1     string
		account2     string
	)

	// read and process lines
	for scanner.Scan() {
		text := scanner.Text()

		// extract the first cluster
		clusterRegex := regexp.MustCompile(`clusters:\s*\[([^\s\]]+)`)
		clusterMatch := clusterRegex.FindStringSubmatch(text)
		if len(clusterMatch) > 1 {
			firstCluster = clusterMatch[1]
			fmt.Println("cluster:", firstCluster)
		}

		// extract account numbers
		accountRegex := regexp.MustCompile(`accounts:\s*\[([^\]]+)\]`)
		accountMatch := accountRegex.FindStringSubmatch(text)
		if len(accountMatch) > 1 {
			accountsStr := accountMatch[1]
			accountNumbers := strings.Split(accountsStr, ",")
			for i, numStr := range accountNumbers {
				num, err := strconv.Atoi(strings.TrimSpace(numStr))
				if err == nil {
					fmt.Printf("- account %d: %d\n", i+1, num)
				}
			}

			account1 = strings.TrimSpace(accountNumbers[0])
			account2 = strings.TrimSpace(accountNumbers[1])
		}

		// get client shards
		cluster1, err := m.storage.GetClientShard(account1)
		if err != nil {
			return err.Error()
		}
		cluster2, err := m.storage.GetClientShard(account2)
		if err != nil {
			return err.Error()
		}

		// update clusters
		if firstCluster == cluster1 {
			// get the cluster's services to remove account2 from cluster2
			services := strings.Split(m.dialer.Nodes[fmt.Sprintf("E%s", cluster2)], ":")
			targetBalance := 0
			for _, svc := range services {
				if _, balance, err := m.dialer.Rebalance(svc, account2, 0, false); err != nil {
					return err.Error()
				} else {
					targetBalance = balance
				}
			}

			// add account2 to the cluster1
			services = strings.Split(m.dialer.Nodes[fmt.Sprintf("E%s", cluster1)], ":")
			for _, svc := range services {
				if _, _, err := m.dialer.Rebalance(svc, account2, targetBalance, true); err != nil {
					return err.Error()
				}
			}

			// update shard
			if err := m.storage.UpdateClientShard(account2, cluster1); err != nil {
				return err.Error()
			}
		} else {
			// get the cluster's services to remove account1 from cluster1
			services := strings.Split(m.dialer.Nodes[fmt.Sprintf("E%s", cluster1)], ":")
			targetBalance := 0
			for _, svc := range services {
				if _, balance, err := m.dialer.Rebalance(svc, account1, 0, false); err != nil {
					return err.Error()
				} else {
					targetBalance = balance
				}
			}

			// add account1 to the cluster2
			services = strings.Split(m.dialer.Nodes[fmt.Sprintf("E%s", cluster2)], ":")
			for _, svc := range services {
				if _, _, err := m.dialer.Rebalance(svc, account1, targetBalance, true); err != nil {
					return err.Error()
				}
			}

			// update shard
			if err := m.storage.UpdateClientShard(account1, cluster2); err != nil {
				return err.Error()
			}
		}
	}

	return "rebalanced"
}
