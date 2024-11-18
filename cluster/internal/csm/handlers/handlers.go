package handlers

import "github.com/F24-CSE535/2pc/cluster/internal/storage"

// NewDatabaseHandler returns an instance of database handler.
func NewDatabaseHandler(st *storage.Database) *DatabaseHandler {
	return &DatabaseHandler{
		storage: st,
	}
}
