package storage

import (
	"context"

	"github.com/F24-CSE535/2pc/cluster/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
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

// IsCollectionEmpty returns true if the collection is empty.
func (d *Database) IsCollectionEmpty() (bool, error) {
	// count the number of documents in the collection
	count, err := d.clientsCollection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetClientBalance returns a balance value by accepting a client.
func (d *Database) GetClientBalance(client string) (int, error) {
	// create a filter for the specified cluster
	filter := bson.M{"client": client}

	// decode the response
	var clientInstance models.Client
	err := d.clientsCollection.FindOne(context.TODO(), filter).Decode(&clientInstance)
	if err != nil {
		return 0, err
	}

	return int(clientInstance.Balance), nil
}
