package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a module that uses mongo-driver library to handle MongoDB queries.
type Database struct {
	conn               *mongo.Client
	shardsCollection   *mongo.Collection
	sessionsCollection *mongo.Collection
}

// NewDatabase opens a MySQL connection and returns an instance of
// database struct.
func NewDatabase(uri string, database string) (*Database, error) {
	// open a new connection to MongoDB cluster
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to open a MongoDB connection: %v", err)
	}

	// create a new instance
	instance := Database{
		conn: conn,
	}

	// create pointers to collections
	instance.shardsCollection = conn.Database(database).Collection("shards")
	instance.sessionsCollection = conn.Database(database).Collection("sessions")

	return &instance, nil
}
