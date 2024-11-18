package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
)

// Manager is responsible for creating consensus state machines.
type Manager struct {
	Storage *storage.Database
	Channel chan *packets.Packet
}

// Initialize accepts a number as the processing unit replicas and starts CSMs.
func (m *Manager) Initialize(replicas int) {
	m.Channel = make(chan *packets.Packet)

	for i := 0; i < replicas; i++ {
		csm := ConsensusStateMachine{
			storage: m.Storage,
			channel: m.Channel,
		}

		// start the CSM inside a go-routine
		go func(c *ConsensusStateMachine) {
			c.Start()
		}(&csm)
	}
}
