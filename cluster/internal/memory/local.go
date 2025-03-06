package memory

import "github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

// SharedMemory is a local storage for processes and handlers.
type SharedMemory struct {
	inBlockStatus bool
	leader        string
	nodeName      string
	clusterName   string
	clusterIPs    []string
	iptable       map[string]string

	ballotNumber            *paxos.BallotNumber
	acceptedNum             *paxos.BallotNumber
	acceptedVal             *paxos.AcceptMsg
	acceptedMsgs            []*paxos.AcceptedMsg
	sessionIdsBallotNumbers map[int]*paxos.BallotNumber
}

// NewSharedMemory returns an instance of shared memory.
func NewSharedMemory(leader, nm, cn string, iptable map[string]string) *SharedMemory {
	instance := &SharedMemory{
		inBlockStatus:           false,
		leader:                  leader,
		nodeName:                nm,
		clusterName:             cn,
		iptable:                 iptable,
		ballotNumber:            &paxos.BallotNumber{Sequence: 0, NodeId: nm},
		acceptedNum:             nil,
		acceptedVal:             nil,
		sessionIdsBallotNumbers: make(map[int]*paxos.BallotNumber),
	}

	instance.SetClusterIPs()

	return instance
}
