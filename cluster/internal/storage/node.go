package storage

import (
	"context"
	"time"

	"github.com/F24-CSE535/2pc/cluster/pkg/enums"
	"github.com/F24-CSE535/2pc/cluster/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

// InsertClusterShard gets shard of a cluster and stores them inside clients collection.
func (d *Database) InsertClusterShard(shard []*models.ClientShard) error {
	// convert shard model to interface
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
func (d *Database) IsClientsCollectionEmpty() (bool, error) {
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
	if err := d.clientsCollection.FindOne(context.TODO(), filter).Decode(&clientInstance); err != nil {
		return 0, err
	}

	return int(clientInstance.Balance), nil
}

// InsertClient adds a new client to the shard schema.
func (d *Database) InsertClient(client string, balance int) error {
	_, err := d.clientsCollection.InsertOne(context.TODO(), &models.Client{
		Client:  client,
		Balance: balance,
	})

	return err
}

// DeleteClient removes a client from schema.
func (d *Database) DeleteClient(client string) error {
	// create a filter for the specified cluster
	filter := bson.M{"client": client}

	_, err := d.clientsCollection.DeleteOne(context.TODO(), filter)

	return err
}

// UpdateClientBalance gets a client and new balance to update the balance value.
func (d *Database) UpdateClientBalance(client string, balance int, set bool) error {
	// create a filter for the specified cluster
	filter := bson.M{"client": client}

	// define the update operation
	var update bson.D
	if set {
		update = bson.D{{Key: "$set", Value: bson.D{{Key: "balance", Value: balance}}}}
	} else {
		update = bson.D{{Key: "$inc", Value: bson.D{{Key: "balance", Value: balance}}}}
	}

	// perform the update query
	_, err := d.clientsCollection.UpdateMany(context.TODO(), filter, update)

	return err
}

// InsertWAL adds a new log to the node's logs.
func (d *Database) InsertWAL(log *models.Log) error {
	// insert log
	_, err := d.logsCollection.InsertOne(context.TODO(), log)

	return err
}

// InsertBatchWALs into the database.
func (d *Database) InsertBatchWALs(logs []*models.Log) error {
	// convert log model to interface
	records := make([]interface{}, 0)
	for _, item := range logs {
		records = append(records, item)
	}

	_, err := d.logsCollection.InsertMany(context.TODO(), records)

	return err
}

// GetWALsBySessionId gets a sessionId and returns the logs for that session.
func (d *Database) GetWALsBySessionId(sessionId int) ([]*models.Log, error) {
	// create a filter for the specified log
	filter := bson.M{
		"session_id": sessionId,
		"message":    enums.WALUpdate,
	}

	// find all documents that match the filter
	cursor, err := d.logsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// decode the results into a slice of Logs structs
	var results []*models.Log
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetWALs returns all write-ahead logs.
func (d *Database) GetWALs() ([]*models.Log, error) {
	// create a filter for all
	filter := bson.M{}

	// find all documents that match the filter
	cursor, err := d.logsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// decode the results into a slice of Logs structs
	var results []*models.Log
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetCommittedWALs returns all committed logs.
func (d *Database) GetCommittedWALs(from int) ([]*models.Log, error) {
	// create logs filter
	filter := bson.M{
		"message":                enums.WALCommit,
		"ballot_number_sequence": bson.M{"$gt": from},
	}

	// find all documents that match the filter
	cursor, err := d.logsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// decode the results into a slice of Logs structs
	var results []*models.Log
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetLogsWithCommittedWALs returns the committed transactions.
func (d *Database) GetLogsWithCommittedWALs(from int) ([]*models.Log, error) {
	// create a filet for commit logs that are after the from input
	commitFilter := bson.M{
		"message":                enums.WALCommit,
		"ballot_number_sequence": bson.M{"$gt": from},
	}

	// run commits query
	commitCursor, err := d.logsCollection.Distinct(context.TODO(), "session_id", commitFilter)
	if err != nil {
		return nil, err
	}

	// create logs filter
	filter := bson.M{
		"message":    enums.WALUpdate,
		"session_id": bson.M{"$in": commitCursor},
	}

	// find all documents that match the filter
	cursor, err := d.logsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// decode the results into a slice of Logs structs
	var results []*models.Log
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// InsertLock gets a record and inserts a lock record.
func (d *Database) InsertLock(record string) error {
	_, err := d.locksCollection.InsertOne(context.TODO(), &models.Lock{
		Record:    record,
		DeletedAt: time.Now().String(),
	})

	return err
}
