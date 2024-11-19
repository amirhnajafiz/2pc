package models

// Log is an entity used in our system's write-ahead logging process.
type Log struct {
	Message   string `bson:"message"`
	Record    string `bson:"record"`
	SessionId int    `bson:"session_id"`
	OldValue  int    `bson:"old_value"`
	NewValue  int    `bson:"new_value"`
}
