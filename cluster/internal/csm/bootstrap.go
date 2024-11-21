package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/handlers"
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/lock"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"

	"go.uber.org/zap"
)

// Manager is responsible for fully creating consensus state machines.
type Manager struct {
	Channel           chan *packets.Packet
	DispatcherChannel chan *packets.Packet
	Memory            *memory.SharedMemory
	Storage           *storage.Database
}

// Initialize accepts a number as the number of processing units, then it starts CSMs.
func (m *Manager) Initialize(logr *zap.Logger, replicas int) {
	// the manager input channel
	m.Channel = make(chan *packets.Packet, 10)
	m.DispatcherChannel = make(chan *packets.Packet, 10)

	// create a new dispatcher
	dis := NewDispatcher(m.DispatcherChannel, m.Channel)

	// create database handler
	dbh := handlers.NewDatabaseHandler(
		client.NewClient(m.Memory.GetNodeName()),
		lock.NewManager(),
		logr.Named("csm-db-handler"),
		m.Memory,
		m.Storage,
	)

	// create paxos handler
	pxh := handlers.NewPaxosHandler(
		m.Channel,
		dis.GetNotifyChannel(),
		client.NewClient(m.Memory.GetNodeName()),
		logr.Named("csm-paxos-handler"),
		m.Memory,
	)

	for i := 0; i < replicas; i++ {
		// create a new CSM
		csm := ConsensusStateMachine{
			channel:         m.Channel,
			databaseHandler: dbh,
			paxosHandler:    pxh,
		}

		// start the CSM inside a go-routine
		go func(c *ConsensusStateMachine, index int) {
			logr.Info("consensus state machine is running", zap.Int("replica number", index))
			c.Start()
		}(&csm, i)
	}
}
