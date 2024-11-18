package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/handlers"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
)

// ConsensusStateMachine is a processing unit that captures packets from gRPC level and passes them to handlers.
type ConsensusStateMachine struct {
	databaseHandler *handlers.DatabaseHandler
	channel         chan *packets.Packet
}

// Start method waits on packets on the input channel, and performs a logic based on packet label.
func (c *ConsensusStateMachine) Start() {
	for {
		// get the gRPC messages
		pkt := <-c.channel

		// case on packet label
		switch pkt.Label {
		case packets.PktRequest:
			c.databaseHandler.Request(pkt.Payload.(*database.RequestMsg).GetTransaction())
		case packets.PktPrepare:
			c.databaseHandler.Prepare()
		case packets.PktCommit:
			c.databaseHandler.Commit()
		case packets.PktAbort:
			c.databaseHandler.Abort()
		}
	}
}
