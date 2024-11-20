package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/handlers"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"
)

// ConsensusStateMachine is a processing unit that captures packets from gRPC level and passes them to handlers.
type ConsensusStateMachine struct {
	databaseHandler *handlers.DatabaseHandler
	paxosHandler    *handlers.PaxosHandler
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
			msg := pkt.Payload.(*database.RequestMsg)
			c.databaseHandler.Request(msg.GetReturnAddress(), msg.GetTransaction())
		case packets.PktPrepare:
			c.databaseHandler.Prepare(pkt.Payload.(*database.PrepareMsg))
		case packets.PktCommit:
			c.databaseHandler.Commit(pkt.Payload.(*database.CommitMsg))
		case packets.PktAbort:
			c.databaseHandler.Abort(int(pkt.Payload.(*database.AbortMsg).GetSessionId()))
		case packets.PktPaxosAccept:
			c.paxosHandler.Accept(pkt.Payload.(*paxos.AcceptMsg))
		case packets.PktPaxosAccepted:
			c.paxosHandler.Accepted(pkt.Payload.(*paxos.AcceptedMsg))
		case packets.PktPaxosCommit:
			c.paxosHandler.Commit(pkt.Payload.(*paxos.CommitMsg))
		}
	}
}
