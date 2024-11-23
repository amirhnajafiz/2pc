package models

import (
	"time"

	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

// Session is a holder for live transactions tracing.
type Session struct {
	Sender       string
	Receiver     string
	Amount       int
	Type         string
	Text         string
	Id           int
	Participants []string
	Acks         []*database.AckMsg
	Replys       []*database.ReplyMsg
	StartedAt    time.Time
}
