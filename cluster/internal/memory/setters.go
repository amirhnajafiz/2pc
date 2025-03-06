package memory

import (
	"fmt"
	"strings"

	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"
)

// SetClusterIPs is used to extract the ip addresses of nodes inside a cluster.
func (s *SharedMemory) SetClusterIPs() {
	// split the cluster endpoints by ':'
	parts := strings.Split(s.iptable[fmt.Sprintf("E%s", s.clusterName)], ":")

	// ip list
	list := make([]string, 0)
	for _, key := range parts {
		if key != s.nodeName {
			list = append(list, s.iptable[key])
		}
	}

	// set cluster IPs
	s.clusterIPs = list
}

// SetLeader updates the current leader id.
func (s *SharedMemory) SetLeader(leader string) {
	s.leader = leader
}

// SetPotentialBallotNumber adds a new ballot-number to the potential list.
func (s *SharedMemory) SetSessionIdBallotNumber(sessionId int, bn *paxos.BallotNumber) {
	s.sessionIdsBallotNumbers[sessionId] = bn
}

// SetAcceptedMessages initializes the accepted messages.
func (s *SharedMemory) SetAcceptedMessages() {
	s.acceptedMsgs = make([]*paxos.AcceptedMsg, 0)
}

// AppendAcceptedMessage adds a new accepted message to accepted messages list.
func (s *SharedMemory) AppendAcceptedMessage(msg *paxos.AcceptedMsg) {
	s.acceptedMsgs = append(s.acceptedMsgs, msg)
}

// ResetAcceptedMessages removes all accepted messages.
func (s *SharedMemory) ResetAcceptedMessages() {
	s.acceptedMsgs = nil
}

// IncBallotNumber increases ballot-number sequence number.
func (s *SharedMemory) IncBallotNumber() {
	s.ballotNumber.Sequence++
}

// SetBallotNumber updates the sequence number of the ballot-number.
func (s *SharedMemory) SetBallotNumber(seq int64) {
	s.ballotNumber.Sequence = seq
}

// SetAcceptedNum sets a new value for accepted-num.
func (s *SharedMemory) SetAcceptedNum(bn *paxos.BallotNumber) {
	s.acceptedNum = bn
}

// SetAcceptedVal sets a new value for accepted-val.
func (s *SharedMemory) SetAcceptedVal(av *paxos.AcceptMsg) {
	s.acceptedVal = av
}

// SetBlockStatus updates the block status
func (s *SharedMemory) SetBlockStatus(status bool) {
	s.inBlockStatus = status
}
