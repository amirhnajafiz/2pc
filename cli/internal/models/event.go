package models

// Event is a struct for storing events in MongoDB collection.
type Event struct {
	Cluster   string `bson:"cluster"`
	Operation string `bson:"operation"`
	CreatedAt string `bson:"created_at"`
	IsDone    bool   `bson:"is_done"`
}
