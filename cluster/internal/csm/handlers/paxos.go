package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/utils"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"go.uber.org/zap"
)

// PaxosHandler contains methods to perform paxos consensus protocol logic.
type PaxosHandler struct {
	client  *client.Client
	logger  *zap.Logger
	memory  *memory.SharedMemory
	channel chan *packets.Packet

	acceptedNum  *paxos.BallotNumber
	acceptedVal  *paxos.AcceptMsg
	acceptedMsgs []*paxos.AcceptedMsg
}

// Request accepts a database request and converts it to paxos request.
func (p *PaxosHandler) Request(req *database.RequestMsg) {
	// create a list for accepted messages
	p.acceptedMsgs = make([]*paxos.AcceptedMsg, 0)

	// increament ballot-number
	p.acceptedNum.Sequence++

	// create paxos request
	msg := paxos.AcceptMsg{
		Request:      utils.ConvertDatabaseRequestToPaxosRequest(req),
		BallotNumber: p.acceptedNum,
		NodeId:       p.memory.GetNodeName(),
		CrossShard:   false,
	}

	// send accept messages
	for _, address := range p.memory.GetClusterIPs() {
		if err := p.client.Accept(address, &msg); err != nil {
			p.logger.Warn("failed to send accept message", zap.Error(err))
		}
	}

	// save the accepted val
	p.acceptedVal = &msg
}

// Prepare accepts a database prepare and converts it to paxos request.
func (p *PaxosHandler) Prepare(req *database.PrepareMsg) {
	// create a list for accepted messages
	p.acceptedMsgs = make([]*paxos.AcceptedMsg, 0)

	// increament ballot-number
	p.acceptedNum.Sequence++

	// create paxos request
	msg := paxos.AcceptMsg{
		Request:      utils.ConvertDatabasePrepareToPaxosRequest(req),
		BallotNumber: p.acceptedNum,
		NodeId:       p.memory.GetNodeName(),
		CrossShard:   true,
	}

	// send accept messages
	for _, address := range p.memory.GetClusterIPs() {
		if err := p.client.Accept(address, &msg); err != nil {
			p.logger.Warn("failed to send accept message", zap.Error(err))
		}
	}

	// save the accepted val
	p.acceptedVal = &msg
}

// Accept gets a new accept message and updates it's datastore and returns an accepted message.
func (p *PaxosHandler) Accept(msg *paxos.AcceptMsg) {
	// update accepted number and accepted value
	p.acceptedNum = msg.GetBallotNumber()
	p.acceptedVal = msg

	// send accepted message
	if err := p.client.Accepted(p.memory.GetFromIPTable(msg.GetNodeId()), p.acceptedNum, p.acceptedVal); err != nil {
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
	if len(p.acceptedMsgs) == 1 {
		for _, address := range p.memory.GetClusterIPs() {
			if err := p.client.Commit(address, p.acceptedNum, p.acceptedVal); err != nil {
				p.logger.Warn("failed to send commit message")
			}
		}

		p.acceptedMsgs = nil
	}

	// send a new request to our own channel
	pkt := packets.Packet{}

	// check for the request type
	if p.acceptedVal.CrossShard {
		pkt.Label = packets.PktDatabasePrepare
		pkt.Payload = &database.PrepareMsg{
			Transaction:   utils.ConvertPaxosRequestToDatabaseTransaction(p.acceptedVal.GetRequest()),
			Client:        p.acceptedVal.Request.GetClient(),
			ReturnAddress: p.acceptedVal.Request.GetReturnAddress(),
		}
	} else {
		pkt.Label = packets.PktDatabaseRequest
		pkt.Payload = &database.RequestMsg{
			Transaction:   utils.ConvertPaxosRequestToDatabaseTransaction(p.acceptedVal.GetRequest()),
			ReturnAddress: p.acceptedVal.Request.GetReturnAddress(),
		}
	}

	p.channel <- &pkt
}

// Commit gets a commit message and creates a new request into the system.
func (p *PaxosHandler) Commit(msg *paxos.CommitMsg) {
	// send a new request to our own channel
	pkt := packets.Packet{}

	// check for the request type
	if p.acceptedVal.CrossShard {
		pkt.Label = packets.PktDatabasePrepare
		pkt.Payload = &database.PrepareMsg{
			Transaction:   utils.ConvertPaxosRequestToDatabaseTransaction(p.acceptedVal.GetRequest()),
			Client:        p.acceptedVal.Request.GetClient(),
			ReturnAddress: p.acceptedVal.Request.GetReturnAddress(),
		}
	} else {
		pkt.Label = packets.PktDatabaseRequest
		pkt.Payload = &database.RequestMsg{
			Transaction:   utils.ConvertPaxosRequestToDatabaseTransaction(p.acceptedVal.GetRequest()),
			ReturnAddress: p.acceptedVal.Request.GetReturnAddress(),
		}
	}

	p.channel <- &pkt
}
