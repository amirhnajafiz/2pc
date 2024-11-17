package models

// Shard is a model that holds each shard's clients.
type Shard struct {
	StartId int
	EndId   int
	Cluster string
	Name    string
	Clients map[string]int
}
