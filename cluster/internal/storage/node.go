package storage

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/pkg/models"
)

// InsertClusterShard gets shard of a cluster and stores them inside clients collection.
func (d *Database) InsertClusterShard(shard []*models.ClientShard) error {
	records := make([]interface{}, 0)
	for _, item := range shard {
		records = append(records, &models.Client{
			Client:  item.Client,
			Balance: item.InitBalance,
		})
	}

	_, err := d.clientsCollection.InsertMany(context.TODO(), records)

	return err
}
