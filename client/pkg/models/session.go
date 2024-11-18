package models

import "github.com/F24-CSE535/2pc/client/pkg/rpc/database"

// Session is a holder for live transactions tracing.
type Session struct {
	Type         string
	Text         string
	Id           int
	Participants int
	Acks         []*database.AckMsg
}
