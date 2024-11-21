package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"go.uber.org/zap"
)

// PaxosHandler contains methods to perform paxos consensus protocol logic.
type PaxosHandler struct {
	client *client.Client
	logger *zap.Logger

	nodeName    string
	clusterName string
	iptable     map[string]string

	channel chan *packets.Packet

	acceptedNum  *paxos.BallotNumber
	acceptedVal  *paxos.AcceptMsg
	acceptedMsgs []*paxos.AcceptedMsg
}

// InitConsensus starts the consensus protocol.
func (p *PaxosHandler) Request() {
	// create a list for accepted messages
	p.acceptedMsgs = make([]*paxos.AcceptedMsg, 0)

	// send accept messages
	if err := p.client.Accept("", nil); err != nil {
		p.logger.Warn("failed to send accept message")
	}
}

// Accept gets a new accept message and updates it's datastore and returns an accepted message.
func (p *PaxosHandler) Accept(msg *paxos.AcceptMsg) {
	// update accepted number and accepted value
	p.acceptedNum = msg.GetBallotNumber()
	p.acceptedVal = msg

	// send accepted
	if err := p.client.Accepted(msg.GetNodeId(), p.acceptedNum, p.acceptedVal); err != nil {
		p.logger.Warn("failed to send accepted message", zap.String("to", msg.GetNodeId()))
	}
}

// Accepted gets a new accepted message and follows the paxos protocol.
func (p *PaxosHandler) Accepted(msg *paxos.AcceptedMsg) {
	// check the consensus is on going
	if p.acceptedMsgs == nil {
		return
	}

	// store the accepted message
	p.acceptedMsgs = append(p.acceptedMsgs, msg)

	// count the messages, if we got the majority send commit messages
	if len(p.acceptedMsgs) == 2 {
		if err := p.client.Commit(); err != nil {
			p.logger.Warn("failed to send commit message")
		}

		p.acceptedMsgs = nil
	}

	// send a new request to our own channel
}

// Commit gets a commit message and creates a new request into the system.
func (p *PaxosHandler) Commit(msg *paxos.CommitMsg) {
	// send a new request to our own channel
}
