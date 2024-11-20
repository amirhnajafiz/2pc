package manager

import (
	"fmt"
	"strconv"

	"github.com/F24-CSE535/2pc/client/pkg/models"
	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

func (m *Manager) Performance() string {
	return fmt.Sprintf("throughput: %f tps, latency: %d ms", m.throughput, m.latency)
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

	// make RPC call
	if balance, err := m.dialer.PrintBalance(cluster, argv[0]); err != nil {
		return fmt.Errorf("server failed: %v", err).Error()
	} else {
		return fmt.Sprintf("%s : %d", argv[0], balance)
	}
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

	// create sessions records
	records := make([]string, 0)
	for _, id := range list {
		if value, ok := m.cache[id]; ok {
			records = append(
				records,
				fmt.Sprintf(
					"transaction %d (%s, %s, %d)",
					value.Id,
					value.Sender,
					value.Receiver,
					value.Amount),
			)
		}
	}

	return records, ""
}

func (m *Manager) Transaction(argc int, argv []string) string {
	// check the number of arguments
	if argc < 3 {
		return "not enough arguments"
	}

	// extract data from input command
	sender := argv[0]
	receiver := argv[1]
	amount, _ := strconv.Atoi(argv[2])
	sessionId := m.memory.GetSession()

	// get shards
	senderCluster, err := m.storage.GetClientShard(sender)
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
	}
	receiverCluster, err := m.storage.GetClientShard(receiver)
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
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
			return fmt.Errorf("server failed: %v", err).Error()
		}
	} else {
		session.Type = "cross-shard"
		session.Participants = []string{senderCluster, receiverCluster}
		session.Acks = make([]*database.AckMsg, 0)

		// for cross-shard send prepare messages to both clusters
		if err := m.dialer.Prepare(senderCluster, sender, sender, receiver, amount, sessionId); err != nil {
			return fmt.Errorf("sender server failed: %v", err).Error()
		}
		if err := m.dialer.Prepare(receiverCluster, receiver, sender, receiver, amount, sessionId); err != nil {
			return fmt.Errorf("receiver server failed: %v", err).Error()
		}
	}

	// save the transaction into cache
	m.cache[sessionId] = &session

	return fmt.Sprintf("transaction %d (%s, %s, %d): sent", sessionId, sender, receiver, amount)
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
