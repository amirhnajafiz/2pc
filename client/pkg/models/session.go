package models

import (
	"time"

	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

// Session is a holder for live transactions tracing.
type Session struct {
	Sender       string               `bson:"sender"`
	Receiver     string               `bson:"receiver"`
	Amount       int                  `bson:"amount"`
	Type         string               `bson:"type"`
	Text         string               `bson:"text"`
	Id           int                  `bson:"id"`
	Participants []string             `bson:"participants"`
	Acks         []*database.AckMsg   `bson:"-"`
	Replys       []*database.ReplyMsg `bson:"-"`
	StartedAt    time.Time            `bson:"-"`
}
