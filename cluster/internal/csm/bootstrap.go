package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/handlers"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"

	"go.uber.org/zap"
)

// Manager is responsible for fully creating consensus state machines.
type Manager struct {
	Storage *storage.Database
	Channel chan *packets.Packet
}

// Initialize accepts a number as the number of processing units, then it starts CSMs.
func (m *Manager) Initialize(logr *zap.Logger, replicas int) {
	// the manager input channel
	m.Channel = make(chan *packets.Packet)

	for i := 0; i < replicas; i++ {
		// create a new CSM
		csm := ConsensusStateMachine{
			databaseHandler: handlers.NewDatabaseHandler(m.Storage, logr.Named("csm-handler")),
			channel:         m.Channel,
		}

		// start the CSM inside a go-routine
		go func(c *ConsensusStateMachine) {
			logr.Info("Consensus State Machine is running")
			c.Start()
		}(&csm)
	}
}
