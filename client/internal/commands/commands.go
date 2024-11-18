package commands

import (
	"fmt"

	"github.com/F24-CSE535/2pc/client/internal/grpc"
	"github.com/F24-CSE535/2pc/client/internal/storage"
)

// Commands is a struct that handles client input commands.
type Commands struct {
	Dialer  *grpc.Dialer
	Storage *storage.Database
}

func (c Commands) PrintBalance(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// get the shard
	cluster, err := c.Storage.GetClientShard(argv[0])
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
	}

	// make RPC call
	if balance, err := c.Dialer.PrintBalance(cluster, argv[0]); err != nil {
		return fmt.Errorf("server failed: %v", err).Error()
	} else {
		return fmt.Sprintf("%s : %d", argv[0], balance)
	}
}
