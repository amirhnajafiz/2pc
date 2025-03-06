package models

// AggregationResult is the output of sharding rebalance.
type AggregationResult struct {
	Account1         string   `bson:"account1"`
	Account2         string   `bson:"account2"`
	TransactionCount int      `bson:"transactionCount"`
	TotalAmount      float64  `bson:"totalAmount"`
	Types            []string `bson:"types"`
	Participants     []string `bson:"participants"`
}
