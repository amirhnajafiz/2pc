package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/F24-CSE535/2pc/cli/internal/database"
	"github.com/F24-CSE535/2pc/cli/internal/utils"
	"github.com/F24-CSE535/2pc/cli/pkg/models"
)

// ShardHandler reads the input datastore and shards files and
// creates the shardings.
type ShardHandler struct{}

func (c *ShardHandler) GetName() string {
	return "shard"
}

func (c *ShardHandler) Execute(argc int, args []string) error {
	// four arguments are needed
	if argc != 3 {
		return fmt.Errorf("mismatch input arguments: count %d expected 3", argc)
	}

	// open database connection
	db, err := database.NewDatabase(args[1], args[2])
	if err != nil {
		return fmt.Errorf("open database failed: %v", err)
	}

	// parse input csv files (args[0] is the shards path)
	shards, err := utils.CSVParseShardsFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse shards: %v", err)
	}

	// create a list of shards
	list := make([]*models.Shard, 0)

	// loop over shards to build our list
	for _, shard := range shards {
		tmp := models.Shard{
			Name:    strings.TrimSpace(shard[0]),
			Cluster: strings.TrimSpace(shard[1]),
			Clients: make(map[string]int),
		}

		// split range by '-'
		parts := strings.Split(shard[2], "-")
		if tmp.StartId, err = strconv.Atoi(strings.TrimSpace(parts[0])); err != nil {
			return fmt.Errorf("failed to parse range %s: %v", shard[0], err)
		}
		if tmp.EndId, err = strconv.Atoi(strings.TrimSpace(parts[1])); err != nil {
			return fmt.Errorf("failed to parse range %s: %v", shard[0], err)
		}

		// create clients
		for i := tmp.StartId; i <= tmp.EndId; i++ {
			tmp.Clients[fmt.Sprintf("%d", i)] = 10
		}

		list = append(list, &tmp)
	}

	// print some metadata
	for _, shard := range list {
		if err := db.InsertShards(shard.DTOClients()); err != nil {
			return fmt.Errorf("failed to inserts %s to database: %v", shard.Name, err)
		}

		fmt.Printf(
			"shard %s mapped to cluster %s with %d clients (%d-%d).\n",
			shard.Name,
			shard.Cluster,
			len(shard.Clients),
			shard.StartId,
			shard.EndId,
		)
	}

	return nil
}
