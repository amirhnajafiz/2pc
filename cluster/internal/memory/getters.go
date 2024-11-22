package memory

import "github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

// GetLeader returns the current leader.
func (s *SharedMemory) GetLeader() string {
	return s.leader
}

// GetNodeName returns the current node.
func (s *SharedMemory) GetNodeName() string {
	return s.nodeName
}

// GetClusterName returns the cluster name.
func (s *SharedMemory) GetClusterName() string {
	return s.clusterName
}

// GetClusterIPs returns the node IPs that are in this cluster.
func (s *SharedMemory) GetClusterIPs() []string {
	return s.clusterIPs
}

// GetFromIPTable returns an address from iptable.
func (s *SharedMemory) GetFromIPTable(key string) string {
	return s.iptable[key]
}

// GetLastCommittedBallotNumber returns the last committed ballot-number.
func (s *SharedMemory) GetLastCommittedBallotNumber() *paxos.BallotNumber {
	return s.lastCommittedBallotNumber
}

// GetBallotNumber returns the node's ballot-number.
func (s *SharedMemory) GetBallotNumber() *paxos.BallotNumber {
	return s.ballotNumber
}

// GetAcceptedNum returns the node's accepted-number.
func (s *SharedMemory) GetAcceptedNum() *paxos.BallotNumber {
	return s.acceptedNum
}

// GetAcceptedVal returns the node's accepted-value.
func (s *SharedMemory) GetAcceptedVal() *paxos.AcceptMsg {
	return s.acceptedVal
}

// IsAcceptedMessagesEmpty returns true if accepted messages is nil.
func (s *SharedMemory) IsAcceptedMessagesEmpty() bool {
	return s.acceptedMsgs == nil
}

// AcceptedMessagesSize returns the len of accepted messages.
func (s *SharedMemory) AcceptedMessagesSize() int {
	return len(s.acceptedMsgs)
}
