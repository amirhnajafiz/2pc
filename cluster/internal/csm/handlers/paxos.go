package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/csm/timers"
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/internal/utils"
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
func (p *PaxosHandler) Request(payload interface{}, isCrossShared bool) {
	if p.memory.GetBlockStatus() {
		return
	}

	// create a list for accepted messages
	p.memory.SetAcceptedMessages()

	// increament ballot-number
	p.memory.IncBallotNumber()

	// create paxos accept message
	msg := paxos.AcceptMsg{
		BallotNumber: p.memory.GetBallotNumber(),
		NodeId:       p.memory.GetNodeName(),
		CrossShard:   isCrossShared,
	}

	// set the request based on prepare type
	if isCrossShared {
		msg.Request = utils.ConvertDatabasePrepareToPaxosRequest(payload.(*database.PrepareMsg))
	} else {
		msg.Request = utils.ConvertDatabaseRequestToPaxosRequest(payload.(*database.RequestMsg))
	}

	// send accept messages
	for _, address := range p.memory.GetClusterIPs() {
		if err := p.client.Accept(address, &msg); err != nil {
			p.logger.Warn("failed to send accept message", zap.Error(err))
		}
	}

	// start consensus timer
	go p.paxosTimer.StartConsensusTimer(msg.GetRequest().GetReturnAddress(), int(msg.GetRequest().GetSessionId()))

	// save the accepted-num and accepted-val
	p.memory.SetAcceptedNum(p.memory.GetBallotNumber())
	p.memory.SetAcceptedVal(&msg)
}

// Accept gets a new accept message and updates it's datastore and returns an accepted message.
func (p *PaxosHandler) Accept(msg *paxos.AcceptMsg) {
	if p.memory.GetBlockStatus() {
		return
	}

	// don't accept old ballot-numbers
	if msg.GetBallotNumber().GetSequence() < p.memory.GetBallotNumber().GetSequence() {
		return
	}

	// don't accept old accepted nums
	if val := p.memory.GetAcceptedNum(); val != nil && val.GetSequence() >= msg.GetBallotNumber().GetSequence() {
		return
	}

	// update accepted number and accepted value if there is nothing stored
	acceptedVal := p.memory.GetAcceptedVal()
	if acceptedVal == nil {
		p.memory.SetAcceptedNum(msg.GetBallotNumber())
		p.memory.SetAcceptedVal(msg)
	}

	// send accepted message back to the leader
	if err := p.client.Accepted(p.memory.GetFromIPTable(msg.GetNodeId()), p.memory.GetAcceptedNum(), acceptedVal); err != nil {
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
	if p.memory.GetAcceptedMessagesSize() < p.majorityAcceptedMessages {
		return
	}

	// stop the consensus timer
	p.paxosTimer.FinishConsensusTimer()

	// select the accepted number and accepted value
	var (
		tmpAcceptedNum *paxos.BallotNumber
		acceptedVal    *paxos.AcceptMsg = p.memory.GetAcceptedVal()
	)
	for _, item := range p.memory.GetAcceptedMessages() {
		if item.GetAcceptedValue() != nil {
			if tmpAcceptedNum == nil || tmpAcceptedNum.GetSequence() < item.GetAcceptedNumber().GetSequence() {
				tmpAcceptedNum = item.GetAcceptedNumber()
				acceptedVal = item.GetAcceptedValue()
			}
		}
	}

	// send commit messages to all servers
	for _, address := range p.memory.GetClusterIPs() {
		if err := p.client.Commit(address, p.memory.GetAcceptedNum(), acceptedVal); err != nil {
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
	if p.memory.GetBlockStatus() {
		return
	}

	// send a new request to our own channel
	pkt := packets.Packet{}

	// get accepted valu from message
	acceptedVal := msg.GetAcceptedValue()

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

	// save the ballot-number in memory map
	p.memory.SetSessionIdBallotNumber(int(msg.GetAcceptedValue().GetRequest().GetSessionId()), msg.GetAcceptedNumber())
	p.memory.SetBallotNumber(p.memory.GetAcceptedNum().GetSequence())

	// reset accepted-num and accepted-val
	p.memory.SetAcceptedNum(nil)
	p.memory.SetAcceptedVal(nil)

	p.csmsChan <- &pkt
}

// Ping gets a ping message, and accepts it if the leader is better.
func (p *PaxosHandler) Ping(msg *paxos.PingMsg) {
	if p.memory.GetBlockStatus() {
		return
	}

	// leader is not good enough (S1 is better than S2)
	if msg.GetNodeId() > p.memory.GetNodeName() {
		return
	}

	// reset the timer
	p.memory.SetLeader(msg.GetNodeId())
	p.leaderTimer.StartLeaderTimer()
	p.leaderTimer.StopLeaderPinger()

	// check the last committed message
	diff := p.memory.GetBallotNumber().GetSequence() - msg.GetLastCommitted().GetSequence()
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
		if err := p.client.Pong(p.memory.GetFromIPTable(msg.GetNodeId()), p.memory.GetBallotNumber()); err != nil {
			p.logger.Warn("failed to send pong message", zap.Error(err), zap.String("to", msg.GetNodeId()))
		}
	}
}

// Pong gets a pong message and syncs the follower.
func (p *PaxosHandler) Pong(msg *paxos.PongMsg) {
	if p.memory.GetBlockStatus() {
		return
	}

	// get paxos items
	pis, err := p.storage.GetLogsWithCommittedWALs(int(msg.GetLastCommitted().GetSequence()))
	if err != nil {
		p.logger.Warn("failed to get paxos items", zap.Error(err))
	}

	// create items
	items := make([]*paxos.SyncItem, 0)
	for _, pi := range pis {
		items = append(items, &paxos.SyncItem{
			Record: pi.Record,
			Value:  int64(pi.NewValue),
		})
	}

	// sync the follower by calling sync
	if err := p.client.Sync(p.memory.GetFromIPTable(msg.GetNodeId()), p.memory.GetBallotNumber(), items); err != nil {
		p.logger.Warn("failed to send sync message", zap.Error(err))
	}
}

// Sync gets a sync message and syncs the node.
func (p *PaxosHandler) Sync(msg *paxos.SyncMsg) {
	if p.memory.GetBlockStatus() {
		return
	}

	// drop the old sync messages
	if msg.GetLastCommitted().GetSequence() < p.memory.GetBallotNumber().GetSequence() {
		return
	}

	// in a loop, update the clients
	for _, item := range msg.GetItems() {
		if err := p.storage.UpdateClientBalance(item.GetRecord(), int(item.GetValue()), false); err != nil {
			p.logger.Warn("failed to update client in sync", zap.Error(err), zap.String("client", item.GetRecord()))
		}
	}

	// update our ballot-number
	p.memory.SetBallotNumber(msg.GetLastCommitted().GetSequence())

	// reset accepted-num and accepted-val
	p.memory.SetAcceptedNum(nil)
	p.memory.SetAcceptedVal(nil)
}

// Block stops all processes in CSMs.
func (p *PaxosHandler) Block() {
	if p.memory.GetBlockStatus() {
		return
	}

	p.memory.SetBlockStatus(true)
	p.leaderTimer.StopLeaderPinger()
	p.leaderTimer.StopLeaderTimer()
}

// Unblock restarts all processes in CSMs.
func (p *PaxosHandler) Unblock() {
	if !p.memory.GetBlockStatus() {
		return
	}

	p.memory.SetBlockStatus(false)
	p.leaderTimer.StartLeaderTimer()
	p.leaderTimer.StartLeaderPinger()
}
