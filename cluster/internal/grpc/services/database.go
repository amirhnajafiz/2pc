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

	Storage           *storage.Database
	Channel           chan *packets.Packet
	DispatcherChannel chan *packets.Packet
}

// Request RPC sends a new request packet to CSMs.
func (d *DatabaseService) Request(_ context.Context, msg *database.RequestMsg) (*emptypb.Empty, error) {
	d.DispatcherChannel <- &packets.Packet{Label: packets.PktPaxosRequest, Payload: msg}

	return &emptypb.Empty{}, nil
}

// Prepare RPC sends a new prepare packet to CSMs.
func (d *DatabaseService) Prepare(_ context.Context, msg *database.PrepareMsg) (*emptypb.Empty, error) {
	d.DispatcherChannel <- &packets.Packet{Label: packets.PktPaxosPrepare, Payload: msg}

	return &emptypb.Empty{}, nil
}

// Commit RPC sends a new commit packet to CSMs.
func (d *DatabaseService) Commit(_ context.Context, msg *database.CommitMsg) (*emptypb.Empty, error) {
	d.Channel <- &packets.Packet{Label: packets.PktDatabaseCommit, Payload: msg}

	return &emptypb.Empty{}, nil
}

// Abort RPC sends a new abort packet to CSMs.
func (d *DatabaseService) Abort(_ context.Context, msg *database.AbortMsg) (*emptypb.Empty, error) {
	d.Channel <- &packets.Packet{Label: packets.PktDatabaseAbort, Payload: msg}

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

// PrintLogs returns all WALs inside this node.
func (d *DatabaseService) PrintLogs(_ *emptypb.Empty, stream database.Database_PrintLogsServer) error {
	wals, err := d.Storage.GetWALs()
	if err != nil {
		return fmt.Errorf("database failed: %v", err)
	}

	// send logs one by one
	for _, wal := range wals {
		if err := stream.Send(&database.LogRsp{
			Record:    wal.Record,
			Message:   wal.Message,
			SessionId: int64(wal.SessionId),
			NewValue:  int64(wal.NewValue),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PrintDatastore returns all committed transactions inside this node.
func (d *DatabaseService) PrintDatastore(_ *emptypb.Empty, stream database.Database_PrintDatastoreServer) error {
	wals, err := d.Storage.GetCommittedWALs(-1)
	if err != nil {
		return fmt.Errorf("database failed: %v", err)
	}

	// send datastore one by one
	for _, wal := range wals {
		if err := stream.Send(&database.DatastoreRsp{
			Record:               wal.Record,
			SessionId:            int64(wal.SessionId),
			BallotNumberSequence: int64(wal.BallotNumberSequence),
			BallotNumberPid:      wal.BallotNumberPid,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Block sends a block packet to the consensus.
func (d *DatabaseService) Block(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	d.Channel <- &packets.Packet{Label: packets.PktDatabaseBlock}

	return &emptypb.Empty{}, nil
}

// Unblock sends a unblock packet to the consensus.
func (d *DatabaseService) Unblock(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	d.Channel <- &packets.Packet{Label: packets.PktDatabaseUnblock}

	return &emptypb.Empty{}, nil
}

// Rebalance updates the node's schema.
func (d *DatabaseService) Rebalance(_ context.Context, msg *database.RebalanceMsg) (*database.RebalanceRsp, error) {
	// add the record if true
	if msg.GetAdd() {
		if err := d.Storage.InsertClient(msg.GetRecord(), int(msg.GetBalance())); err != nil {
			return nil, fmt.Errorf("database failed: %v", err)
		}
	} else {
		// get the client balance, remove the client and return the response
		balance, err := d.Storage.GetClientBalance(msg.GetRecord())
		if err != nil {
			return nil, fmt.Errorf("database failed: %v", err)
		}

		// remove the client
		if err := d.Storage.DeleteClient(msg.GetRecord()); err != nil {
			return nil, fmt.Errorf("database failed: %v", err)
		}

		return &database.RebalanceRsp{
			Record:  msg.GetRecord(),
			Balance: int64(balance),
		}, nil
	}

	return &database.RebalanceRsp{}, nil
}
