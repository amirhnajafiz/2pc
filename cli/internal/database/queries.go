package database

import (
	"context"
	"time"

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

// InsertEvent creates a new event in the database.
func (d *Database) InsertEvent(event *models.Event) error {
	event.CreatedAt = time.Now().String()
	event.IsDone = false

	_, err := d.eventsCollection.InsertOne(context.TODO(), event)

	return err
}
