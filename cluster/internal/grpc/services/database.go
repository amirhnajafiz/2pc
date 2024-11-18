package services

import (
	"context"
	"fmt"

	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
	"github.com/F24-CSE535/2pc/cluster/pkg/rpc/database"

	"google.golang.org/protobuf/types/known/emptypb"
)

// DatabaseService is used for handling database RPCs.
type DatabaseService struct {
	database.UnimplementedDatabaseServer

	Storage *storage.Database
	Channel chan *packets.Packet
}

// Request RPC sends a new packet to CSMs.
func (d *DatabaseService) Request(_ context.Context, msg *database.RequestMsg) (*emptypb.Empty, error) {
	d.Channel <- &packets.Packet{Label: packets.PktRequest, Payload: msg}

	return &emptypb.Empty{}, nil
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
