package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/timers"
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/internal/utils"
	"github.com/F24-CSE535/2pc/cluster/pkg/models"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"go.uber.org/zap"
)

// PaxosHandler contains methods to perform paxos consensus protocol logic.
type PaxosHandler struct {
	client      *client.Client
	logger      *zap.Logger
	memory      *memory.SharedMemory
	storage     *storage.Database
	leaderTimer *timers.LeaderTimer
	paxosTimer  *timers.PaxosTimer

	majorityAcceptedMessages int

	dispatcherNotifyChan chan bool
	csmsChan             chan *packets.Packet
}

// Request accepts a database request and converts it to paxos request.
func (p *PaxosHandler) Request(req *database.RequestMsg) {
	// create a list for accepted messages
	p.memory.SetAcceptedMessages()

	// increament ballot-number
	p.memory.IncBallotNumber()
	p.memory.SetAcceptedNum(p.memory.GetBallotNumber())

	// create paxos request
	msg := paxos.AcceptMsg{
		Request:      utils.ConvertDatabaseRequestToPaxosRequest(req),
		BallotNumber: p.memory.GetAcceptedNum(),
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
	go p.paxosTimer.StartConsensusTimer(req.GetReturnAddress(), int(req.GetTransaction().GetSessionId()))

	// save the accepted val
	p.memory.SetAcceptedVal(&msg)
}

// Prepare accepts a database prepare and converts it to paxos request.
func (p *PaxosHandler) Prepare(req *database.PrepareMsg) {
	// create a list for accepted messages
	p.memory.SetAcceptedMessages()

	// increament ballot-number
	p.memory.IncBallotNumber()
	p.memory.SetAcceptedNum(p.memory.GetBallotNumber())

	// create paxos request
	msg := paxos.AcceptMsg{
		Request:      utils.ConvertDatabasePrepareToPaxosRequest(req),
		BallotNumber: p.memory.GetAcceptedNum(),
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
	go p.paxosTimer.StartConsensusTimer(req.GetReturnAddress(), int(req.GetTransaction().GetSessionId()))

	// save the accepted val
	p.memory.SetAcceptedVal(&msg)
}

// Accept gets a new accept message and updates it's datastore and returns an accepted message.
func (p *PaxosHandler) Accept(msg *paxos.AcceptMsg) {
	// don't accept old ballot-numbers
	if msg.GetBallotNumber().GetSequence() < p.memory.GetBallotNumber().GetSequence() {
		return
	}

	// update accepted number and accepted value
	p.memory.SetAcceptedNum(msg.GetBallotNumber())
	p.memory.SetAcceptedVal(msg)
	p.memory.SetBallotNumber(msg.GetBallotNumber().GetSequence())

	// send accepted message
	if err := p.client.Accepted(p.memory.GetFromIPTable(msg.GetNodeId()), p.memory.GetAcceptedNum(), p.memory.GetAcceptedVal()); err != nil {
		p.logger.Warn("failed to send accepted message", zap.String("to", msg.GetNodeId()))
	}
}

// Accepted gets a new accepted message and follows the paxos protocol.
func (p *PaxosHandler) Accepted(msg *paxos.AcceptedMsg) {
	// check the consensus is on going
	if p.memory.IsAcceptedMessagesEmpty() {
		return
	}

	// store the accepted message
	p.memory.AppendAcceptedMessage(msg)

	// count the messages, if we got the majority
	if p.memory.AcceptedMessagesSize() < p.majorityAcceptedMessages {
		return
	}

	// stop the consensus timer
	p.paxosTimer.FinishConsensusTimer()

	// send commit messages
	for _, address := range p.memory.GetClusterIPs() {
		if err := p.client.Commit(address, p.memory.GetAcceptedNum(), p.memory.GetAcceptedVal()); err != nil {
			p.logger.Warn("failed to send commit message")
		}
	}

	// reset all accepted messages
	p.memory.ResetAcceptedMessages()

	// send a new commit message to our own channel
	pkt := packets.Packet{
		Label: packets.PktPaxosCommit,
		Payload: &paxos.CommitMsg{
			AcceptedNumber: p.memory.GetAcceptedNum(),
			AcceptedValue:  p.memory.GetAcceptedVal(),
		},
	}

	// send the commit to our own channel and notify the dispatcher
	p.csmsChan <- &pkt
	p.dispatcherNotifyChan <- true
}

// Commit gets a commit message and creates a new request into the system.
func (p *PaxosHandler) Commit(msg *paxos.CommitMsg) {
	// send a new request to our own channel
	pkt := packets.Packet{}

	// get accepted val
	acceptedVal := p.memory.GetAcceptedVal()

	// check for the request type
	if acceptedVal.CrossShard {
		pkt.Label = packets.PktDatabasePrepare
		pkt.Payload = &database.PrepareMsg{
			Transaction:   utils.ConvertPaxosRequestToDatabaseTransaction(acceptedVal.GetRequest()),
			Client:        acceptedVal.Request.GetClient(),
			ReturnAddress: acceptedVal.Request.GetReturnAddress(),
		}
	} else {
		pkt.Label = packets.PktDatabaseRequest
		pkt.Payload = &database.RequestMsg{
			Transaction:   utils.ConvertPaxosRequestToDatabaseTransaction(acceptedVal.GetRequest()),
			ReturnAddress: acceptedVal.Request.GetReturnAddress(),
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

	p.csmsChan <- &pkt
}

// Ping gets a ping message, and accepts it if the leader is better.
func (p *PaxosHandler) Ping(msg *paxos.PingMsg) {
	// leader is not good enough (S1 is better than S2)
	if msg.GetNodeId() > p.memory.GetNodeName() {
		return
	}

	// reset the timer
	p.memory.SetLeader(msg.GetNodeId())
	p.leaderTimer.StartLeaderTimer()
	p.leaderTimer.StopLeaderPinger()

	// check the last committed message
	diff := p.memory.GetLastCommittedBallotNumber().GetSequence() - msg.GetLastCommitted().GetSequence()
	if diff > 0 {
		// sync the leader by generating a pong message
		p.csmsChan <- &packets.Packet{
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

	// update our last committed and ballot-number
	p.memory.ResetLastCommittedBallotNumber(msg.GetLastCommitted())
	p.memory.SetBallotNumber(msg.GetLastCommitted().GetSequence() + 1)
}
