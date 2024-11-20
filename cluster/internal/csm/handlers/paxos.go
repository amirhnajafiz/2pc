package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"go.uber.org/zap"
)

// PaxosHandler contains methods to perform paxos consensus protocol logic.
type PaxosHandler struct {
	client *client.Client
	logger *zap.Logger

	outputChannel chan bool

	acceptedNum  *paxos.BallotNumber
	acceptedVal  *paxos.AcceptMsg
	acceptedMsgs []*paxos.AcceptedMsg
}

// GetOutputChannel returns the paxos handler ourput channel.
func (p *PaxosHandler) GetOutputChannel() chan bool {
	return p.outputChannel
}

// InitConsensus starts the consensus protocol.
func (p *PaxosHandler) InitConsensus() {
	p.outputChannel = make(chan bool)
	p.acceptedMsgs = make([]*paxos.AcceptedMsg, 0)
}

// EndConsensus ends the consensus protocol.
func (p *PaxosHandler) EndConsensus() {
	p.outputChannel = nil
	p.acceptedMsgs = nil
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
	if p.outputChannel == nil {
		return
	}

	// store the accepted message
	p.acceptedMsgs = append(p.acceptedMsgs, msg)

	// count the messages
	if len(p.acceptedMsgs) == 2 {
		p.outputChannel <- true
	}
}

// Commit gets a commit message and creates a new request into the system.
func (p *PaxosHandler) Commit(msg *paxos.CommitMsg) {

}
