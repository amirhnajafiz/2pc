package models

// PaxosItem is a combinator of accepted_num and accepted_val that we store in MongoDB.
type PaxosItem struct {
	BallotNumberNum int    `bson:"ballot_number_num"`
	BallotNumberPid string `bson:"ballot_number_pid"`
	Client          string `bson:"client"`
	Sender          string `bson:"sender"`
	Receiver        string `bson:"receiver"`
	Amount          int    `bson:"amount"`
	SessionId       int    `bson:"session_id"`
	IsCommitted     bool   `bson:"is_committed"`
}
