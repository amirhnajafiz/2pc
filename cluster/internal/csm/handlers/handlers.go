package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/lock"
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"go.uber.org/zap"
)

// NewDatabaseHandler returns an instance of database handler.
func NewDatabaseHandler(
	st *storage.Database,
	mem *memory.SharedMemory,
	logr *zap.Logger,
	client *client.Client,
	lm *lock.Manager,
) *DatabaseHandler {
	return &DatabaseHandler{
		memory:  mem,
		logger:  logr,
		storage: st,
		manager: lm,
		client:  client,
	}
}

// NewPaxosHandler returns an instance paxos handler.
func NewPaxosHandler(
	channel chan *packets.Packet,
	logr *zap.Logger,
	client *client.Client,
	mem *memory.SharedMemory,
) *PaxosHandler {
	return &PaxosHandler{
		memory:      mem,
		channel:     channel,
		logger:      logr,
		client:      client,
		acceptedNum: &paxos.BallotNumber{Sequence: 0, NodeId: mem.GetNodeName()},
		acceptedVal: nil,
	}
}
