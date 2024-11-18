package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
)

// ConsensusStateMachine is a processing unit that captures packets and calls handlers.
type ConsensusStateMachine struct {
	storage *storage.Database
	channel chan *packets.Packet
}

func (c *ConsensusStateMachine) Start() {
	for {
		// get the gRPC messages
		pkt := <-c.channel

		switch pkt.Label {
		case packets.PktRequest:
			break
		case packets.PktPrepare:
			break
		case packets.PktCommit:
			break
		case packets.PktAbort:
			break
		}
	}
}
