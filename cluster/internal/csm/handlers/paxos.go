package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"go.uber.org/zap"
)

// PaxosHandler contains methods to perform paxos consensus protocol logic.
type PaxosHandler struct {
	client *client.Client
	logger *zap.Logger
}
