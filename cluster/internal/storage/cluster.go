package storage

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

// GetClusterShard returns all items from global database that are belong to this cluster.
func (d *Database) GetClusterShard() ([]*models.ClientShard, error) {
	// create a filter for the specified cluster
	filter := bson.M{"cluster": d.cluster}

	// find all documents that match the filter
	cursor, err := d.shardsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// decode the results into a slice of ClientShard structs
	var results []*models.ClientShard
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetEvents returns a list of the system events.
func (d *Database) GetEvents() ([]*models.Event, error) {
	// create a filter for the specified events
	filter := bson.M{
		"cluster": d.cluster,
		"is_done": false,
	}

	// find all documents that match the filter
	cursor, err := d.eventsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// decode the results into a slice of Events structs
	var results []*models.Event
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (d *Database) UpdateEvents() error {
	// create a filter for the specified events
	filter := bson.M{
		"cluster": d.cluster,
		"is_done": false,
	}

	// define the update operation
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_done", Value: true}}}}

	// perform the update
	_, err := d.eventsCollection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
