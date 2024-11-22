package memory

import "github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

// SharedMemory is a local storage for processes and handlers.
type SharedMemory struct {
	leader      string
	nodeName    string
	clusterName string
	clusterIPs  []string
	iptable     map[string]string

	potentialCommittedBallotNumbers map[int]*paxos.BallotNumber
	lastCommittedBallotNumber       *paxos.BallotNumber
}

// NewSharedMemory returns an instance of shared memory.
func NewSharedMemory(leader, nm, cn string, iptable map[string]string) *SharedMemory {
	instance := &SharedMemory{
		leader:                          leader,
		nodeName:                        nm,
		clusterName:                     cn,
		iptable:                         iptable,
		potentialCommittedBallotNumbers: make(map[int]*paxos.BallotNumber),
		lastCommittedBallotNumber:       nil,
	}

	instance.SetClusterIPs()

	return instance
}
