package models

// Lock is a wrapper for locks in MongoDB.
type Lock struct {
	Record    string `bson:"record"`
	DeletedAt string `bson:"deleted_at"`
}
