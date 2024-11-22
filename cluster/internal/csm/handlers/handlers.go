package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/timers"
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/lock"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"

	"go.uber.org/zap"
)

// NewDatabaseHandler returns an instance of database handler.
func NewDatabaseHandler(
	client *client.Client,
	lm *lock.Manager,
	logr *zap.Logger,
	mem *memory.SharedMemory,
	st *storage.Database,
) *DatabaseHandler {
	return &DatabaseHandler{
		memory:  mem,
		logger:  logr,
		storage: st,
		manager: lm,
		client:  client,
	}
}

// NewPaxosHandler returns an instance paxos handler.
func NewPaxosHandler(
	channel chan *packets.Packet,
	channelNotify chan bool,
	client *client.Client,
	logr *zap.Logger,
	mem *memory.SharedMemory,
	st *storage.Database,
) *PaxosHandler {
	instance := &PaxosHandler{
		memory:               mem,
		storage:              st,
		logger:               logr,
		client:               client,
		csmsChan:             channel,
		dispatcherNotifyChan: channelNotify,
		leaderTimerChan:      make(chan bool),
		leaderPingChan:       make(chan bool),
		consensusTimerChan:   make(chan bool),
	}

	// create timers
	lt := timers.NewLeaderTimer(client, logr.Named("leader-timer"), mem, instance.leaderPingChan, instance.leaderTimerChan)
	instance.paxosTimer = timers.NewPaxosTimer(client, logr.Named("paxos-timer"), mem, instance.consensusTimerChan, instance.dispatcherNotifyChan)

	// start the leader timer and leader pinger
	go lt.LeaderTimer()
	go lt.LeaderPinger()

	return instance
}
