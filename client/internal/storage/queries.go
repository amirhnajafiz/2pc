package storage

import (
	"context"
	"fmt"

	"github.com/F24-CSE535/2pc/client/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetClientShard gets a client id and returns its shard data.
func (d *Database) GetClientShard(client string) (string, error) {
	// create a filter for the specified cluster
	filter := bson.M{"client": client}
	findOptions := options.FindOne().SetProjection(bson.M{"cluster": 1, "_id": 0})

	// decode the response
	var shard struct {
		Cluster string `bson:"cluster"`
	}
	err := d.shardsCollection.FindOne(context.TODO(), filter, findOptions).Decode(&shard)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no shard found for client: %s", client)
		}

		return "", fmt.Errorf("error fetching shard: %v", err)
	}

	return shard.Cluster, nil
}

// InsertSession stores a session for future system rebalance.
func (d *Database) InsertSession(session *models.Session) error {
	// create a new item to store
	var item struct {
		Sender   string   `bson:"sender"`
		Receiver string   `bson:"receiver"`
		Shards   []string `bson:"shards"`
	}
	item.Sender = session.Sender
	item.Receiver = session.Receiver
	item.Shards = session.Participants

	_, err := d.crossShardsCollection.InsertOne(context.TODO(), &item)

	return err
}
