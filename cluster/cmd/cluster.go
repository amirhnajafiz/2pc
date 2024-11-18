package cmd

import (
	"fmt"
	"time"

	"github.com/F24-CSE535/2pc/cluster/internal/config"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/pkg/logger"

	"go.uber.org/zap"
)

// Cluster is a manager that monitors events and performs operations on a cluster nodes.
type Cluster struct {
	ConfigPath  string
	ClusterName string

	cfg      config.Config
	database *storage.Database

	ports int
	nodes []chan bool
}

func (c *Cluster) Main() error {
	// load cluster configs
	c.cfg = config.New(c.ConfigPath)
	c.ports = c.cfg.Subnet

	// create new console zap logger and a new file logger
	logr := logger.NewCoreLogger(c.cfg.LogLevel)
	floger := logger.NewFileLogger(c.cfg.LogLevel)

	// open global database connection
	gdb, err := storage.NewClusterDatabase(c.cfg.MongoDB, c.cfg.Database, c.ClusterName)
	if err != nil {
		return fmt.Errorf("failed to open global database connection: %v", err)
	}

	// assign the global database to cluster database field
	c.database = gdb

	// open a list of nodes channels
	c.nodes = make([]chan bool, 0)

	// in a loop, monitor the events
	perioud := time.Duration(c.cfg.WatchInterval) * time.Second
	for {
		time.Sleep(perioud)

		// get events
		events, err := c.database.GetEvents()
		if err != nil {
			logr.Warn("failed to get events", zap.Error(err))
			continue
		}

		// loop over events
		for _, event := range events {
			switch event.Operation {
			case "scale-up":
				c.scaleUp(floger)
				logr.Info("scaled up", zap.Int("nodes", len(c.nodes)))
			}
		}

		// update events
		if err := c.database.UpdateEvents(); err != nil {
			logr.Warn("failed update events", zap.Error(err))
		}
	}
}

// scaleUp creates a new node instance.
func (c *Cluster) scaleUp(loger *zap.Logger) error {
	name := fmt.Sprintf("S%d", len(c.nodes)+1)

	// open the new node database
	ndb, err := storage.NewNodeDatabase(c.cfg.MongoDB, c.ClusterName, name)
	if err != nil {
		return fmt.Errorf("failed to open %s database connection: %v", name, err)
	}

	if isEmpty, err := ndb.IsCollectionEmpty(); err != nil {
		return fmt.Errorf("failed to check %s clients collection status: %v", name, err)
	} else if isEmpty {
		// clone the shards into the node database
		sh, err := c.database.GetClusterShard()
		if err != nil {
			return fmt.Errorf("failed to get global cluster shard: %v", err)
		}
		if err := ndb.InsertClusterShard(sh); err != nil {
			return fmt.Errorf("failed to create %s clients collections: %v", name, err)
		}
	}

	// create a new node
	n := node{
		logger:             loger,
		database:           ndb,
		terminationChannel: make(chan bool),
	}

	// set the node channel
	c.nodes = append(c.nodes, n.terminationChannel)

	// set the node's port
	port := c.ports
	c.ports++

	// start the node
	go n.main(port, name)

	return nil
}
