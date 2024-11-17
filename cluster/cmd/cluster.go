package cmd

import (
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
	gdb, err := storage.NewClusterDatabase(cfg.MongoDB, cfg.Database)
	if err != nil {
		return err
	}

	// open new node database
	ndb, err := storage.NewNodeDatabase(cfg.MongoDB, c.ClusterName, "S1")
	if err != nil {
		return err
	}

	cfg.GRPCPort = 6001

	sh, err := gdb.GetClusterShard()
	if err != nil {
		return err
	}

	if err := ndb.InsertClusterShard(sh); err != nil {
		return err
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
