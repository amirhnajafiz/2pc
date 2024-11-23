package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/config/paxos"
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
	Cfg     *paxos.Config
	Memory  *memory.SharedMemory
	Storage *storage.Database

	Channel           chan *packets.Packet
	DispatcherChannel chan *packets.Packet
}

// Initialize accepts a number as the number of processing units, then it starts CSMs.
func (m *Manager) Initialize(logr *zap.Logger) {
	// the manager input channel
	m.Channel = make(chan *packets.Packet, m.Cfg.CSMBufferSize)
	m.DispatcherChannel = make(chan *packets.Packet, m.Cfg.CSMBufferSize)

	// create a new dispatcher
	dis := NewDispatcher(m.DispatcherChannel, m.Channel, m.Memory)

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
		m.Cfg,
		m.Channel,
		dis.GetNotifyChannel(),
		client.NewClient(m.Memory.GetNodeName()),
		logr.Named("csm-paxos-handler"),
		m.Memory,
		m.Storage,
	)

	for i := 0; i < m.Cfg.CSMReplicas; i++ {
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
