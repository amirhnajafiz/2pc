package memory

import "github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

// SharedMemory is a local storage for processes and handlers.
type SharedMemory struct {
	leader      string
	nodeName    string
	clusterName string
	clusterIPs  []string
	iptable     map[string]string

	ballotNumber *paxos.BallotNumber
	acceptedNum  *paxos.BallotNumber
	acceptedVal  *paxos.AcceptMsg
	acceptedMsgs []*paxos.AcceptedMsg

	lastCommitted    *paxos.BallotNumber
	ballotNumbersMap map[int]*paxos.BallotNumber
}

// NewSharedMemory returns an instance of shared memory.
func NewSharedMemory(leader, nm, cn string, iptable map[string]string) *SharedMemory {
	instance := &SharedMemory{
		leader:           leader,
		nodeName:         nm,
		clusterName:      cn,
		iptable:          iptable,
		ballotNumber:     &paxos.BallotNumber{Sequence: 0, NodeId: nm},
		acceptedNum:      &paxos.BallotNumber{Sequence: 0, NodeId: nm},
		lastCommitted:    &paxos.BallotNumber{Sequence: 0, NodeId: nm},
		acceptedVal:      nil,
		ballotNumbersMap: make(map[int]*paxos.BallotNumber),
	}

	instance.SetClusterIPs()

	return instance
}
