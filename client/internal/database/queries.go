package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetClientShard gets a client id and returns its shard data.
func (d *Database) GetClientShard(client string) (string, error) {
	// create a filter for the specified cluster
	filter := bson.M{"client": client}
	findOptions := options.FindOne().SetProjection(bson.M{"cluster": 1, "_id": 0})

	// find first document that match the filter
	cursor := d.shardsCollection.FindOne(context.TODO(), filter, findOptions)
	if err := cursor.Err(); err != nil {
		return "", err
	}

	// decode the response
	var shard string
	if err := cursor.Decode(&shard); err != nil {
		return "", err
	}

	return shard, nil
}
