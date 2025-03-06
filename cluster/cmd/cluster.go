package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/F24-CSE535/2pc/cluster/internal/config"
	"github.com/F24-CSE535/2pc/cluster/internal/storage"
	"github.com/F24-CSE535/2pc/cluster/internal/utils"
	"github.com/F24-CSE535/2pc/cluster/pkg/logger"

	"go.uber.org/zap"
)

// Cluster is a manager that monitors events and performs operations on a cluster nodes.
type Cluster struct {
	ConfigPath  string
	IPTablePath string

	cfg      config.Config
	database *storage.Database

	activeReplicas int
	ports          int
	iptable        map[string]string

	wg sync.WaitGroup
}

func (c *Cluster) Main() error {
	// load cluster configs
	c.cfg = config.New(c.ConfigPath)
	c.ports = c.cfg.Subnet
	c.activeReplicas = 0

	// load iptables
	nodes, err := utils.IPTableParseFile(c.IPTablePath)
	if err != nil {
		return fmt.Errorf("failed to open iptables: %v", err)
	}
	c.iptable = nodes

	// create a new file logger for our nodes
	logr := logger.NewFileLogger(c.cfg.LogLevel, c.cfg.ReplicasStartingIndex)

	// open global database connection
	db, err := storage.NewClusterDatabase(c.cfg.MongoDB, c.cfg.Database, c.cfg.ClusterName)
	if err != nil {
		return fmt.Errorf("failed to open global database connection: %v", err)
	}
	c.database = db

	// for init replicas create instances
	for i := 0; i < c.cfg.Replicas; i++ {
		if err := c.scaleUp(logr); err != nil {
			log.Printf("failed to start new replica: %v", err)
		}
	}

	// wait for all replicas
	c.wg.Wait()

	return nil
}

// scaleUp creates a new node instance.
func (c *Cluster) scaleUp(loger *zap.Logger) error {
	name := fmt.Sprintf("S%d", c.cfg.ReplicasStartingIndex+c.activeReplicas)

	// select the first node as init leader
	leader := name
	if c.activeReplicas > 0 {
		leader = fmt.Sprintf("S%d", c.cfg.ReplicasStartingIndex)
	}

	// open the new node database
	db, err := storage.NewNodeDatabase(c.cfg.MongoDB, c.cfg.ClusterName, name)
	if err != nil {
		return fmt.Errorf("failed to open %s database connection: %v", name, err)
	}

	// check if the collection is empty
	if isEmpty, err := db.IsClientsCollectionEmpty(); err != nil {
		return fmt.Errorf("failed to check %s clients collection status: %v", name, err)
	} else if isEmpty {
		// clone the shards into the node database
		sh, err := c.database.GetClusterShard()
		if err != nil {
			return fmt.Errorf("failed to get global cluster shard: %v", err)
		}
		if err := db.InsertClusterShard(sh); err != nil {
			return fmt.Errorf("failed to create %s clients collections: %v", name, err)
		}
	}

	// create a new node
	n := node{
		cfg:      &c.cfg.PaxosConfig,
		logger:   loger.Named(name),
		database: db,
		cluster:  c.cfg.ClusterName,
		iptable:  c.iptable,
		leader:   leader,
	}

	// set the node's port
	port := c.ports
	c.ports++

	// start the node
	go n.main(port, name)

	// increase the wait-group
	c.wg.Add(1)

	c.activeReplicas++
	log.Printf("scaled up; current nodes: %d\n", c.activeReplicas)

	return nil
}
