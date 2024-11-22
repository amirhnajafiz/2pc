package timers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"

	"go.uber.org/zap"
)

// NewLeaderTimer returns an instance of leader timer.
func NewLeaderTimer(
	client *client.Client,
	logger *zap.Logger,
	memory *memory.SharedMemory,
	leaderTo,
	leaderPi int,
) *LeaderTimer {
	instance := LeaderTimer{
		leaderTimeout:      leaderTo,
		leaderPingInterval: leaderPi,
		client:             client,
		logger:             logger,
		memory:             memory,
		leaderPingChan:     make(chan bool),
		leaderTimerChan:    make(chan bool),
	}

	// start the leader timer and leader pinger
	go instance.leaderTimer()
	go instance.leaderPinger()

	return &instance
}

// NewPaxosTimer returns an instance of paxos timer.
func NewPaxosTimer(
	cto int,
	client *client.Client,
	logger *zap.Logger,
	memory *memory.SharedMemory,
	dispatcherNotifyChan chan bool,
) *PaxosTimer {
	return &PaxosTimer{
		consensusTimeout:     cto,
		client:               client,
		logger:               logger,
		memory:               memory,
		consensusTimerChan:   make(chan bool),
		dispatcherNotifyChan: dispatcherNotifyChan,
	}
}
