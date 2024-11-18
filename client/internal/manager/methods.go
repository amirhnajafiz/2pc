package manager

import (
	"fmt"
	"strconv"
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

func (m *Manager) Transaction(argc int, argv []string) string {
	// check the number of arguments
	if argc < 4 {
		return "not enough arguments"
	}

	// extract data from input command
	sender := argv[0]
	receiver := argv[1]
	amount, _ := strconv.Atoi(argv[2])
	sessionId := 0

	// get shards
	senderCluster, err := m.storage.GetClientShard(sender)
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
	}
	receiverCluster, err := m.storage.GetClientShard(receiver)
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
	}

	// check for inter or cross shard
	if senderCluster == receiverCluster {
		if err := m.dialer.Request(senderCluster, sender, receiver, amount, sessionId); err != nil {
			return fmt.Errorf("server failed: %v", err).Error()
		}
	}

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
