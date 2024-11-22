package timers

import (
	"time"

	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/pkg/enums"

	"go.uber.org/zap"
)

// PaxosTimer will be used in paxos handler to start requests' timer.
type PaxosTimer struct {
	client *client.Client
	logger *zap.Logger
	memory *memory.SharedMemory

	consensusTimeout time.Duration

	consensusTimerChan   chan bool
	dispatcherNotifyChan chan bool
}

// consensusTimers starts a consensus timer and if it hits timeout it will generate a timeout message.
func (p *PaxosTimer) StartConsensusTimer(ra string, sessionId int) {
	// create a new timer and start it
	timer := time.NewTimer(p.consensusTimeout)

	select {
	case <-p.consensusTimerChan:
		timer.Stop()
	case <-timer.C:
		// close everything and reply with timeout
		p.memory.ResetAcceptedMessages()
		p.logger.Info("consensus timeout", zap.Int("session id", sessionId))

		if err := p.client.Reply(ra, enums.RespConsensusFailed, sessionId); err != nil {
			p.logger.Warn("failed to send reply message", zap.Error(err), zap.String("to", ra))
		}

		// accept next request
		p.dispatcherNotifyChan <- true
	}
}

// FinishConsensusTimer sends a notify signal to the consensus timer channel.
func (p *PaxosTimer) FinishConsensusTimer() {
	p.consensusTimerChan <- true
}
