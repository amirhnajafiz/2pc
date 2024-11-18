package cmd

import (
	"fmt"
	"log"
	"sync"
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
	Index       int

	cfg      config.Config
	database *storage.Database

	wg    sync.WaitGroup
	ports int
	nodes []chan bool
}

func (c *Cluster) Main() error {
	// load cluster configs
	c.cfg = config.New(c.ConfigPath)
	c.ports = c.cfg.Subnet

	// create a new file logger for our nodes
	logr := logger.NewFileLogger(c.cfg.LogLevel)

	// open global database connection
	db, err := storage.NewClusterDatabase(c.cfg.MongoDB, c.cfg.Database, c.ClusterName)
	if err != nil {
		return fmt.Errorf("failed to open global database connection: %v", err)
	}
	c.database = db

	// make list of nodes channels
	c.nodes = make([]chan bool, 0)

	// for init replicas create instances
	for i := 0; i < c.cfg.Replicas; i++ {
		if err := c.scaleUp(logr); err != nil {
			log.Printf("failed to start new replica: %v", err)
		}
	}

	// if the watch interval was greater than zero, in a loop monitor events
	if c.cfg.WatchInterval > 0 {
		perioud := time.Duration(c.cfg.WatchInterval) * time.Second
		for {
			time.Sleep(perioud)

			// get all not done events
			events, err := c.database.GetEvents()
			if err != nil {
				log.Printf("failed to get events: %v\n", err)
				continue
			}

			// loop over events
			for _, event := range events {
				switch event.Operation {
				case "scale-up":
					if err := c.scaleUp(logr); err != nil {
						log.Printf("failed to scale-up: %v", err)
					}
				case "scale-down":
					c.scaleDown()
				}
			}

			// update events
			if err := c.database.UpdateEvents(); err != nil {
				log.Printf("failed to update events: %v\n", err)
			}
		}
	}

	// wait for all replicas
	c.wg.Wait()

	return nil
}

// scaleUp creates a new node instance.
func (c *Cluster) scaleUp(loger *zap.Logger) error {
	name := fmt.Sprintf("S%d", c.Index+len(c.nodes))

	// open the new node database
	db, err := storage.NewNodeDatabase(c.cfg.MongoDB, c.ClusterName, name)
	if err != nil {
		return fmt.Errorf("failed to open %s database connection: %v", name, err)
	}

	// check if the collection is empty
	if isEmpty, err := db.IsCollectionEmpty(); err != nil {
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
		logger:             loger.Named(name),
		database:           db,
		terminationChannel: make(chan bool),
	}

	// set the node channel
	c.nodes = append(c.nodes, n.terminationChannel)

	// set the node's port
	port := c.ports
	c.ports++

	// start the node
	go n.main(port)

	// increase the wait-group
	c.wg.Add(1)

	log.Printf("scaled up; current nodes: %d\n", len(c.nodes))

	return nil
}

// scaleDown removes the last node from the nodes list.
func (c *Cluster) scaleDown() {
	last := len(c.nodes) - 1

	// terminate the process
	c.nodes[last] <- true
	c.nodes = c.nodes[:last]

	// decrese the wait-group
	c.wg.Done()
	c.ports--

	log.Printf("scaled down; current nodes: %d\n", len(c.nodes))
}
