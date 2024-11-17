package cmd

import (
	"fmt"

	"github.com/F24-CSE535/2pc/cluster/internal/config"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/logger"
)

type Cluster struct {
	ConfigPath  string
	ClusterName string
}

func (c Cluster) Main() error {
	// load configs
	cfg := config.New(c.ConfigPath)

	// create new zap logger
	logr := logger.NewLogger(cfg.LogLevel)

	// open global database connection
	gdb, err := storage.NewClusterDatabase(cfg.MongoDB, cfg.Database, c.ClusterName)
	if err != nil {
		return fmt.Errorf("failed to open global database connection: %v", err)
	}

	// open new node database
	ndb, err := storage.NewNodeDatabase(cfg.MongoDB, c.ClusterName, "S1")
	if err != nil {
		return fmt.Errorf("failed to open node database connection: %v", err)
	}

	cfg.GRPCPort = 6001

	sh, err := gdb.GetClusterShard()
	if err != nil {
		return fmt.Errorf("failed to get cluster shard: %v", err)
	}

	if err := ndb.InsertClusterShard(sh); err != nil {
		return fmt.Errorf("failed to create new node collections: %v", err)
	}

	// create a new node
	n := node{
		cfg:                cfg,
		logger:             logr.Named("S1"),
		database:           ndb,
		terminationChannel: make(chan bool),
	}

	n.main()

	return nil
}
