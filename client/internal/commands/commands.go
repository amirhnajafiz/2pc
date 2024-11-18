package commands

import (
	"fmt"

	"github.com/F24-CSE535/2pc/client/internal/grpc"
	"github.com/F24-CSE535/2pc/client/internal/storage"
)

// Commands is a struct that handles client input commands.
type Commands struct {
	dialer  *grpc.Dialer
	storage *storage.Database

	throughput float64
	latency    int64
	count      int
}

// NewCommands returns a new commands instance.
func NewCommands(dialer *grpc.Dialer, storage *storage.Database) *Commands {
	return &Commands{
		dialer:     dialer,
		storage:    storage,
		throughput: 0,
		latency:    0,
		count:      0,
	}
}

func (c *Commands) Performance() string {
	return fmt.Sprintf("throughput: %f tps, latency: %d ms", c.throughput, c.latency)
}

func (c *Commands) PrintBalance(argc int, argv []string) string {
	// check the number of arguments
	if argc < 1 {
		return "not enough arguments"
	}

	// get the shard
	cluster, err := c.storage.GetClientShard(argv[0])
	if err != nil {
		return fmt.Errorf("database failed: %v", err).Error()
	}

	// make RPC call
	if balance, err := c.dialer.PrintBalance(cluster, argv[0]); err != nil {
		return fmt.Errorf("server failed: %v", err).Error()
	} else {
		return fmt.Sprintf("%s : %d", argv[0], balance)
	}
}
