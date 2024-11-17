package models

// Shard is a model that holds each shard's clients.
type Shard struct {
	StartId int
	EndId   int
	Cluster string
	Name    string
	Clients map[string]int
}

// DTOClients returns a list of client shards from a given shard.
func (s Shard) DTOClients() []*ClientShard {
	list := make([]*ClientShard, len(s.Clients))

	index := 0
	for key, value := range s.Clients {
		list[index].Client = key
		list[index].InitBalance = value
		list[index].Cluster = s.Cluster
		list[index].Shard = s.Name

		index++
	}

	return list
}
