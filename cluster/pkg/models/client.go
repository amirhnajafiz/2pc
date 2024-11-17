package models

// Client is a key-value pair that will be stored inside MongoDB.
type Client struct {
	Client  string `bson:"client"`
	Balance int    `bson:"balance"`
}
