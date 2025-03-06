package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/handlers"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"
)

// ConsensusStateMachine is a processing unit that captures packets from gRPC level and passes them to handlers.
type ConsensusStateMachine struct {
	channel         chan *packets.Packet
	databaseHandler *handlers.DatabaseHandler
	paxosHandler    *handlers.PaxosHandler
}

// Start method waits on packets on the input channel, and performs a logic based on packet label.
func (c *ConsensusStateMachine) Start() {
	for {
		// get the gRPC messages
		pkt := <-c.channel

		// case on packet label
		switch pkt.Label {
		case packets.PktDatabaseRequest: // comes from the paxos handler
			c.databaseHandler.Request(pkt.Payload.(*database.RequestMsg))
		case packets.PktDatabasePrepare: // comes from the paxos handler
			c.databaseHandler.Prepare(pkt.Payload.(*database.PrepareMsg))
		case packets.PktDatabaseCommit: // comes from the gRPC
			c.databaseHandler.Commit(pkt.Payload.(*database.CommitMsg))
		case packets.PktDatabaseAbort: // comes from the gRPC
			c.databaseHandler.Abort(int(pkt.Payload.(*database.AbortMsg).GetSessionId()))
		case packets.PktPaxosRequest: // comes from the gRPC
			c.paxosHandler.Request(pkt.Payload.(*database.RequestMsg), false)
		case packets.PktPaxosPrepare: // comes from the gRPC
			c.paxosHandler.Request(pkt.Payload.(*database.PrepareMsg), true)
		case packets.PktPaxosAccept: // comes from the gRPC
			c.paxosHandler.Accept(pkt.Payload.(*paxos.AcceptMsg))
		case packets.PktPaxosAccepted: // comes from the gRPC
			c.paxosHandler.Accepted(pkt.Payload.(*paxos.AcceptedMsg))
		case packets.PktPaxosCommit: // comes from the gRPC
			c.paxosHandler.Commit(pkt.Payload.(*paxos.CommitMsg))
		case packets.PktPaxosPing: // comes from the gRPC
			c.paxosHandler.Ping(pkt.Payload.(*paxos.PingMsg))
		case packets.PktPaxosPong: // comes from both gRPC and paxos handler
			c.paxosHandler.Pong(pkt.Payload.(*paxos.PongMsg))
		case packets.PktPaxosSync: // comes from the gRPC
			c.paxosHandler.Sync(pkt.Payload.(*paxos.SyncMsg))
		case packets.PktDatabaseBlock: // comes from the gRPC
			c.paxosHandler.Block()
		case packets.PktDatabaseUnblock: // comes from the gRPC
			c.paxosHandler.Unblock()
		}
	}
}
