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

	acceptedNum *paxos.BallotNumber
	acceptedVal *paxos.AcceptMsg
}

// Accept get's a new accept message and updates it's datastore and returns an accepted message.
func (p *PaxosHandler) Accept(msg *paxos.AcceptMsg) {
	// update accepted number and accepted value
	p.acceptedNum = msg.GetBallotNumber()
	p.acceptedVal = msg

	// send accepted
	if err := p.client.Accepted(msg.GetNodeId(), p.acceptedNum, p.acceptedVal); err != nil {
		p.logger.Warn("failed to send accepted message", zap.String("to", msg.GetNodeId()))
	}
}

func (p *PaxosHandler) Accepted(msg *paxos.AcceptedMsg) {

}

func (p *PaxosHandler) Commit(msg *paxos.CommitMsg) {

}
