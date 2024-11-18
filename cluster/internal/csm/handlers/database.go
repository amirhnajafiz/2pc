package handlers

import "github.com/F24-CSE535/2pc/cluster/internal/storage"

// DatabaseHandler contains methods to perform database logic.
type DatabaseHandler struct {
	storage *storage.Database
}

func (d DatabaseHandler) Request() {

}

func (d DatabaseHandler) Prepare() {

}

func (d DatabaseHandler) Commit() {

}

func (d DatabaseHandler) Abort() {

}
