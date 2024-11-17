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
