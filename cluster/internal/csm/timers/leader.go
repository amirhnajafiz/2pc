package timers

import (
	"time"

	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"

	"go.uber.org/zap"
)

// LeaderTimer is a struct for managing leader timeout and leader ping operations.
type LeaderTimer struct {
	client *client.Client
	logger *zap.Logger
	memory *memory.SharedMemory

	leaderTimeout      time.Duration
	leaderPingInterval time.Duration

	leaderPingChan  chan bool
	leaderTimerChan chan bool
}

// StartLeaderTimer sends a true signal to leader timer channel.
func (p *LeaderTimer) StartLeaderTimer() {
	p.leaderTimerChan <- true
}

// StopLeaderTimer sends a false signal to leader timer channel.
func (p *LeaderTimer) StopLeaderTimer() {
	p.leaderTimerChan <- false
}

// leader timer is a go-routine that waits on packets from the leader.
// if it does not get enough responses in time, it will create a leader timeout packet.
func (p *LeaderTimer) leaderTimer() {
	// create a new timer and start it
	timer := time.NewTimer(p.leaderTimeout)

	// leader timer while-loop
	for {
		// stop the timer if we are leader
		if p.memory.GetLeader() == p.memory.GetNodeName() {
			timer.Stop()
		}

		select {
		case value := <-p.leaderTimerChan:
			if value {
				p.logger.Debug("accepting new leader", zap.String("current leader", p.memory.GetLeader()))
				timer.Reset(p.leaderTimeout)
			} else {
				timer.Stop()
			}
		case <-timer.C:
			// the node itself becomes the leader
			p.logger.Debug("leader timeout", zap.String("current leader", p.memory.GetLeader()))
			p.memory.SetLeader(p.memory.GetNodeName())
			p.StartLeaderPinger()
		}
	}
}

// StartLeaderPinger sends a true signal to leader pinger channel.
func (p *LeaderTimer) StartLeaderPinger() {
	p.leaderPingChan <- true
}

// StopLeaderPinger sends a false signal to leader pinger channel.
func (p *LeaderTimer) StopLeaderPinger() {
	p.leaderPingChan <- false
}

// leaderPinger starts pinging other servers until it gets stop by a better leader.
func (p *LeaderTimer) leaderPinger() {
	// create a new timer and start it
	timer := time.NewTimer(p.leaderPingInterval)

	// leader pinger while-loop
	for {
		// stop the timer if we are not leader
		if p.memory.GetLeader() != p.memory.GetNodeName() {
			timer.Stop()
		}

		select {
		case value := <-p.leaderPingChan:
			if value {
				timer.Reset(p.leaderPingInterval)
			} else {
				timer.Stop()
			}
		case <-timer.C:
			// send a ping request to everyone
			for _, address := range p.memory.GetClusterIPs() {
				if err := p.client.Ping(address, p.memory.GetBallotNumber()); err != nil {
					p.logger.Warn("failed to send ping message", zap.Error(err), zap.String("to", address))
				}
			}

			timer.Reset(p.leaderPingInterval)
		}
	}
}
