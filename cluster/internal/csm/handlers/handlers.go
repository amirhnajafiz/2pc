package handlers

import (
	"github.com/F24-CSE535/2pc/cluster/internal/grpc/client"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"go.uber.org/zap"
)

// NewDatabaseHandler returns an instance of database handler.
func NewDatabaseHandler(st *storage.Database, logr *zap.Logger, client *client.Client) *DatabaseHandler {
	return &DatabaseHandler{
		logger:  logr,
		storage: st,
		client:  client,
	}
}
