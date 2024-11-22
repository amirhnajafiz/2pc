package handlers

import (
	"time"

	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/internal/utils"
	"github.com/F24-CSE535/2pc/cluster/pkg/enums"
	"github.com/F24-CSE535/2pc/cluster/pkg/models"
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
	storage *storage.Database

	channel chan *packets.Packet

	notify    chan bool
	timer     chan bool
	leader    chan bool
	consensus chan bool

	acceptedNum  *paxos.BallotNumber
	acceptedVal  *paxos.AcceptMsg
	acceptedMsgs []*paxos.AcceptedMsg
}

// leader timer is a go-routine that waits on packets from the leader.
// if it does not get enough responses in time, it will create a leader timeout packet.
func (p *PaxosHandler) leaderTimer() {
	// create a new timer and start it
	timer := time.NewTimer(2 * time.Second)

	// leader timer while-loop
	for {
		// stop the timer if we are leader
		if p.memory.GetLeader() == p.memory.GetNodeName() {
			timer.Stop()
		}

		select {
		case value := <-p.timer:
			if value {
				p.logger.Debug("accepting new leader", zap.String("current leader", p.memory.GetLeader()))
				timer.Reset(2 * time.Second)
			} else {
				timer.Stop()
			}
		case <-timer.C:
			// the node itself becomes the leader
			p.logger.Debug("leader timeout", zap.String("current leader", p.memory.GetLeader()))
			p.memory.SetLeader(p.memory.GetNodeName())
			p.leader <- true
		}
	}
}

// leaderPinger starts pinging other servers until it gets stop by a better leader.
func (p *PaxosHandler) leaderPinger() {
	// create a new timer and start it
	timer := time.NewTimer(5 * time.Second)

	// leader pinger while-loop
	for {
		// stop the timer if we are not leader
		if p.memory.GetLeader() != p.memory.GetNodeName() {
			timer.Stop()
		}

		select {
		case value := <-p.leader:
			if value {
				timer.Reset(5 * time.Second)
			} else {
				timer.Stop()
			}
		case <-timer.C:
			// send a ping request to everyone
			for _, address := range p.memory.GetClusterIPs() {
				if err := p.client.Ping(address, p.memory.GetLastCommittedBallotNumber()); err != nil {
					p.logger.Warn("failed to send ping message", zap.Error(err), zap.String("to", address))
				}
			}

			timer.Reset(5 * time.Second)
		}
	}
}

// consensusTimers starts a consensus timer and if it hits timeout it will generate a timeout message.
func (p *PaxosHandler) consensusTimer(ra string, sessionId int) {
	// create a new timer and start it
	timer := time.NewTimer(5 * time.Second)

	select {
	case <-p.consensus:
		timer.Stop()
	case <-timer.C:
		// close everything and reply with timeout
		p.acceptedMsgs = nil
		p.logger.Info("consensus timeout", zap.Int("session id", sessionId))

		if err := p.client.Reply(ra, enums.RespConsensusFailed, sessionId); err != nil {
			p.logger.Warn("failed to send reply message", zap.Error(err), zap.String("to", ra))
		}
	}
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

	// start consensus timer
	go p.consensusTimer(req.GetReturnAddress(), int(req.GetTransaction().GetSessionId()))

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

	// start consensus timer
	go p.consensusTimer(req.GetReturnAddress(), int(req.GetTransaction().GetSessionId()))

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

	// count the messages, if we got the majority
	if len(p.acceptedMsgs) < 1 {
		return
	}

	// stop the consensus timer
	p.consensus <- true

	// send commit messages
	for _, address := range p.memory.GetClusterIPs() {
		if err := p.client.Commit(address, p.acceptedNum, p.acceptedVal); err != nil {
			p.logger.Warn("failed to send commit message")
		}
	}

	// reset all accepted messages
	p.acceptedMsgs = nil

	// send a new commit message to our own channel
	pkt := packets.Packet{
		Label: packets.PktPaxosCommit,
		Payload: &paxos.CommitMsg{
			AcceptedNumber: p.acceptedNum,
			AcceptedValue:  p.acceptedVal,
		},
	}

	// send the commit to our own channel and notify the dispatcher
	p.channel <- &pkt
	p.notify <- true
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

	// save the paxos item into storage
	if err := p.storage.InsertPaxosItem(&models.PaxosItem{
		BallotNumberNum: int(msg.GetAcceptedNumber().GetSequence()),
		BallotNumberPid: msg.GetAcceptedNumber().GetNodeId(),
		Client:          msg.GetAcceptedValue().GetRequest().GetClient(),
		Sender:          msg.GetAcceptedValue().GetRequest().GetSender(),
		Receiver:        msg.GetAcceptedValue().GetRequest().GetReceiver(),
		Amount:          int(msg.GetAcceptedValue().GetRequest().GetAmount()),
		SessionId:       int(msg.GetAcceptedValue().GetRequest().GetSessionId()),
	}); err != nil {
		p.logger.Warn("failed to store paxos item", zap.Error(err))
	}

	// save the ballot-number in memory
	p.memory.SetPotentialBallotNumber(int(msg.GetAcceptedValue().GetRequest().GetSessionId()), msg.GetAcceptedNumber())

	p.channel <- &pkt
}

// Ping gets a ping message, and accepts it if the leader is better.
func (p *PaxosHandler) Ping(msg *paxos.PingMsg) {
	// leader is not good enough (S1 is better than S2)
	if msg.GetNodeId() > p.memory.GetNodeName() {
		return
	}

	// reset the timer
	p.memory.SetLeader(msg.GetNodeId())
	p.timer <- true
	p.leader <- false

	// check the last committed message
	diff := p.memory.GetLastCommittedBallotNumber().GetSequence() - msg.GetLastCommitted().GetSequence()
	if diff > 0 {
		// sync the leader by generating a pong message
		p.channel <- &packets.Packet{
			Label: packets.PktPaxosPong,
			Payload: &paxos.PongMsg{
				LastCommitted: msg.GetLastCommitted(),
				NodeId:        msg.GetNodeId(),
			},
		}
	} else if diff < 0 {
		// demand a sync by calling pong
		if err := p.client.Pong(p.memory.GetFromIPTable(msg.GetNodeId()), p.memory.GetLastCommittedBallotNumber()); err != nil {
			p.logger.Warn("failed to send pong message", zap.Error(err), zap.String("to", msg.GetNodeId()))
		}
	}
}

// Pong gets a pong message and syncs the follower.
func (p *PaxosHandler) Pong(msg *paxos.PongMsg) {
	// get paxos items
	pis, err := p.storage.GetPaxosItems(int(msg.GetLastCommitted().GetSequence()))
	if err != nil {
		p.logger.Warn("failed to get paxos items", zap.Error(err))
	}

	// create items
	items := make([]*paxos.SyncItem, 0)
	for _, pi := range pis {
		tmp := paxos.SyncItem{
			Record: pi.Client,
		}

		if pi.Client == pi.Sender {
			tmp.Value = int64(-1 * pi.Amount)
		} else {
			tmp.Value = int64(pi.Amount)
		}

		items = append(items, &tmp)
	}

	// sync the follower by calling sync
	if err := p.client.Sync(p.memory.GetFromIPTable(msg.GetNodeId()), p.memory.GetLastCommittedBallotNumber(), items); err != nil {
		p.logger.Warn("failed to send sync message", zap.Error(err))
	}
}

// Sync gets a sync message and syncs the node.
func (p *PaxosHandler) Sync(msg *paxos.SyncMsg) {
	// drop the old sync messages
	if msg.GetLastCommitted().GetSequence() < p.memory.GetLastCommittedBallotNumber().GetSequence() {
		return
	}

	// in a loop, update the clients
	for _, item := range msg.GetItems() {
		if err := p.storage.UpdateClientBalance(item.GetRecord(), int(item.GetValue()), false); err != nil {
			p.logger.Warn("failed to update client in sync", zap.Error(err), zap.String("client", item.GetRecord()))
		}
	}

	// update our last committed
	p.memory.ResetLastCommittedBallotNumber(msg.GetLastCommitted())
}
