package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/lock"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/paxos"

	"go.uber.org/zap"
)

// NewDatabaseHandler returns an instance of database handler.
func NewDatabaseHandler(st *storage.Database, logr *zap.Logger, client *client.Client, lm *lock.Manager) *DatabaseHandler {
	return &DatabaseHandler{
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
	nodeName string,
	nodes []string,
	iptable map[string]string,
) *PaxosHandler {
	return &PaxosHandler{
		channel:     channel,
		logger:      logr,
		client:      client,
		acceptedNum: &paxos.BallotNumber{Sequence: 0, NodeId: nodeName},
		acceptedVal: nil,
		nodeName:    nodeName,
		nodes:       nodes,
		iptable:     iptable,
	}
}
