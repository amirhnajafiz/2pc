package storage

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetClientBalance returns a balance value by accepting a client.
func (d *Database) GetClientBalance(client string) (int, error) {
	// create a filter for the specified cluster
	filter := bson.M{"client": client}
	findOptions := options.FindOne().SetProjection(bson.M{"balance": 1, "_id": 0})

	// find first document that match the filter
	cursor := d.clientsCollection.FindOne(context.TODO(), filter, findOptions)
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	// decode the response
	var balance int
	if err := cursor.Decode(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
