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

// UpdateClientShard accepts a client and shard to update the client shard.
func (d *Database) UpdateClientShard(client, shard string) error {
	// create a filter for the specified cluster
	filter := bson.M{"client": client}

	// define the update operation
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "cluster", Value: shard}}}}

	// perform the update query
	_, err := d.shardsCollection.UpdateMany(context.TODO(), filter, update)

	return err
}

// InsertSession stores a session for future system rebalance.
func (d *Database) InsertSession(session *models.Session) error {
	_, err := d.sessionsCollection.InsertOne(context.TODO(), &session)

	return err
}

// GetSessionById accepts a sessionId and returns the session that is stored in the database.
func (d *Database) GetSessionById(sessionId int) (*models.Session, error) {
	// create a filter fot the specific session
	filter := bson.M{"id": sessionId}

	// find the session
	var session models.Session
	err := d.sessionsCollection.FindOne(context.TODO(), filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no session found for id: %d", sessionId)
		}

		return nil, fmt.Errorf("error fetching session: %v", err)
	}

	return &session, nil
}
