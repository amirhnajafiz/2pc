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
	leaderPingChan chan bool,
	leaderTimerChan chan bool,
) *LeaderTimer {
	return &LeaderTimer{
		client:          client,
		logger:          logger,
		memory:          memory,
		leaderPingChan:  leaderPingChan,
		leaderTimerChan: leaderTimerChan,
	}
}

// NewPaxosTimer returns an instance of paxos timer.
func NewPaxosTimer(
	client *client.Client,
	logger *zap.Logger,
	memory *memory.SharedMemory,
	dispatcherNotifyChan chan bool,
) *PaxosTimer {
	return &PaxosTimer{
		client:               client,
		logger:               logger,
		memory:               memory,
		consensusTimerChan:   make(chan bool),
		dispatcherNotifyChan: dispatcherNotifyChan,
	}
}
