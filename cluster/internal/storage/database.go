package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a module that uses mongo-driver library to handle MongoDB queries.
type Database struct {
	cluster string

	conn              *mongo.Client
	shardsCollection  *mongo.Collection
	eventsCollection  *mongo.Collection
	clientsCollection *mongo.Collection
	logsCollection    *mongo.Collection
	locksCollection   *mongo.Collection
}

// NewClusterDatabase opens a MongoDB connection and returns an instance of database struct.
func NewClusterDatabase(uri string, database string, cluster string) (*Database, error) {
	// open a new connection to MongoDB cluster
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to open a MongoDB connection: %v", err)
	}

	// create a new instance
	instance := Database{
		conn:    conn,
		cluster: cluster,
	}

	// create pointers to collections
	instance.shardsCollection = conn.Database(database).Collection("shards")
	instance.eventsCollection = conn.Database(database).Collection("events")

	return &instance, nil
}

// NewNodeDatabase opens a MongoDB connection and returns an instance of database struct for nodes.
func NewNodeDatabase(uri string, database string, node string) (*Database, error) {
	// open a new connection to MongoDB cluster
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to open a MongoDB connection: %v", err)
	}

	// create a new instance
	instance := Database{
		conn:    conn,
		cluster: database,
	}

	// create pointers to collections
	instance.clientsCollection = conn.Database(database).Collection(fmt.Sprintf("%s_clients", node))
	instance.logsCollection = conn.Database(database).Collection(fmt.Sprintf("%s_logs", node))
	instance.locksCollection = conn.Database(database).Collection(fmt.Sprintf("%s_locks", node))

	return &instance, nil
}
