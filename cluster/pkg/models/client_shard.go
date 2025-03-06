package models

// ClientShard is a key-object pair that will be stored inside MongoDB.
type ClientShard struct {
	Client      string `bson:"client"`
	Shard       string `bson:"shard"`
	Cluster     string `bson:"cluster"`
	InitBalance int    `bson:"init_balance"`
}
