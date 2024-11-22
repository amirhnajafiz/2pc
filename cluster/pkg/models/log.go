package models

// Log is an entity used in our system's write-ahead logging process.
type Log struct {
	Message      string `bson:"message"`
	Record       string `bson:"record"`
	SessionId    int    `bson:"session_id"`
	NewValue     int    `bson:"new_value"`
	BallotNumber int    `bson:"ballot_number"`
}
