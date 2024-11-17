package services

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"
)

// DatabaseService is used for handling database RPCs.
type DatabaseService struct {
	database.UnimplementedDatabaseServer

	Storage storage.Database
}

// PrintBalance accepts a printbalance message and returns a printbalance response.
func (d *DatabaseService) PrintBalance(_ context.Context, msg *database.PrintBalanceMsg) (*database.PrintBalanceRsp, error) {
	return nil, nil
}
