package database

import (
	"context"

	"github.com/F24-CSE535/2pc/cli/internal/models"
)

// InsertShards creates a shards collection in the database.
func (d *Database) InsertShards(list []*models.ClientShard) error {
	records := make([]interface{}, 0)
	for _, item := range list {
		records = append(records, item)
	}

	_, err := d.shardsCollection.InsertMany(context.TODO(), records)

	return err
}
