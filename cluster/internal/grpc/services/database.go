package services

import (
	"context"
	"fmt"

	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"

	"go.uber.org/zap"
)

// DatabaseService is used for handling database RPCs.
type DatabaseService struct {
	database.UnimplementedDatabaseServer

	Logger  *zap.Logger
	Storage *storage.Database
}

// PrintBalance accepts a printbalance message and returns a printbalance response.
func (d *DatabaseService) PrintBalance(_ context.Context, msg *database.PrintBalanceMsg) (*database.PrintBalanceRsp, error) {
	balance, err := d.Storage.GetClientBalance(msg.GetClient())
	if err != nil {
		return nil, fmt.Errorf("database failed: %v", err)
	}

	return &database.PrintBalanceRsp{
		Client:  msg.GetClient(),
		Balance: int64(balance),
	}, nil
}
