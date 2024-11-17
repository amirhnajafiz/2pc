package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/F24-CSE535/2pc/cli/internal/models"
	"github.com/F24-CSE535/2pc/cli/internal/utils"
)

// ShardHandler reads the input datastore and shards files and
// creates the shardings.
type ShardHandler struct{}

func (c *ShardHandler) GetName() string {
	return "shard"
}

func (c *ShardHandler) Execute(argc int, args []string) error {
	// two arguments are needed
	if argc != 2 {
		return fmt.Errorf("mismatch input arguments: count%d expected 2", argc)
	}

	// parse input csv files (args[0] is datastore path and args[1] is shards path)
	datastore, err := utils.CSVParseDatastoreFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse datastore: %v", err)
	}
	shards, err := utils.CSVParseShardsFile(args[1])
	if err != nil {
		return fmt.Errorf("failed to parse shards: %v", err)
	}

	// create a list of shards
	list := make([]*models.Shard, 0)

	// loop over shards to build our list
	for _, shard := range shards {
		tmp := models.Shard{
			Name:    shard[0],
			Cluster: shard[1],
			Clients: make(map[string]int),
		}

		// split range by '-'
		parts := strings.Split(shard[2], "-")
		if tmp.StartId, err = strconv.Atoi(parts[0]); err != nil {
			return fmt.Errorf("failed to parse range %s: %v", shard[0], err)
		}
		if tmp.EndId, err = strconv.Atoi(parts[1]); err != nil {
			return fmt.Errorf("failed to parse range %s: %v", shard[0], err)
		}

		list = append(list, &tmp)
	}

	// loop over datastore to put each item into its shard
	for key, value := range datastore {
		index, _ := strconv.Atoi(strings.Replace(key, "S", "", -1)) // remove the S from client name

		for _, item := range list {
			// client index should be in the range of shard
			if item.StartId >= index && index < item.EndId {
				item.Clients[key] = value
			}
		}
	}

	fmt.Println(list)

	return nil
}
