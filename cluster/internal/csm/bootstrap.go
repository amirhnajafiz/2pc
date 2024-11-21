package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/handlers"
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/lock"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"

	"go.uber.org/zap"
)

// Manager is responsible for fully creating consensus state machines.
type Manager struct {
	NodeName    string
	ClusterName string
	IPTable     map[string]string

	LockManager *lock.Manager
	Storage     *storage.Database
	Channel     chan *packets.Packet
}

// Initialize accepts a number as the number of processing units, then it starts CSMs.
func (m *Manager) Initialize(logr *zap.Logger, replicas int) {
	// the manager input channel
	m.Channel = make(chan *packets.Packet)

	for i := 0; i < replicas; i++ {
		// create a new CSM
		csm := ConsensusStateMachine{
			databaseHandler: handlers.NewDatabaseHandler(m.Storage, logr.Named("csm-db-handler"), client.NewClient(m.NodeName), m.LockManager),
			paxosHandler: handlers.NewPaxosHandler(
				m.Channel, logr.Named("csm-paxos-handler"),
				client.NewClient(m.NodeName),
				m.NodeName,
				m.ClusterName,
				m.IPTable,
			),
			channel: m.Channel,
		}

		// start the CSM inside a go-routine
		go func(c *ConsensusStateMachine, index int) {
			logr.Info("consensus state machine is running", zap.Int("replica number", index))
			c.Start()
		}(&csm, i)
	}
}
