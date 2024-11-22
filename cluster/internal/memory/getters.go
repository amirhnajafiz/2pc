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
