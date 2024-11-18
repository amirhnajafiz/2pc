package manager

import (
	"fmt"

	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

// handleReply accepts a reply message for an active session.
func (m *Manager) handleReply(msg *database.ReplyMsg) {
	if session, ok := m.cache[int(msg.GetSessionId())]; ok {
		session.Text = msg.GetText()
	}
}

// handleAck accepts an ack message for an active session.
func (m *Manager) handleAck(msg *database.AckMsg) {
	if session, ok := m.cache[int(msg.GetSessionId())]; ok {
		// append the ack to the list
		session.Acks = append(session.Acks, msg)

		// check for the number of acks
		if len(session.Acks) == session.Participants {
			// if any is aborted, then abort all
			for _, item := range session.Acks {
				if item.IsAborted {
					fmt.Println("aborted")
				}
			}

			// if all are committed, then commit all
			fmt.Println("send commit")
		}
	}
}
