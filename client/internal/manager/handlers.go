package manager

import (
	"log"
	"time"

	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

// handleReply accepts a reply message for an active session.
func (m *Manager) handleReply(msg *database.ReplyMsg) {
	if session, ok := m.cache[int(msg.GetSessionId())]; ok {
		// append the reply to the list
		session.Replys = append(session.Replys, msg)

		// check for the number of replys
		if len(session.Replys) == len(session.Participants) {
			fn := time.Now()

			// return the message to client
			session.Text = msg.GetText()
			m.output <- session

			// update performance metrics
			du := fn.Sub(session.StartedAt).Nanoseconds() / 1000000
			m.throughput = append(m.throughput, float64(1000/du))
			m.latency = append(m.latency, float64(du))
		}
	}
}

// handleAck accepts an ack message for an active session.
func (m *Manager) handleAck(msg *database.AckMsg) {
	if session, ok := m.cache[int(msg.GetSessionId())]; ok {
		// append the ack to the list
		session.Acks = append(session.Acks, msg)

		// check for the number of acks
		if len(session.Acks) == len(session.Participants) {
			// if any is aborted, then abort all
			for _, item := range session.Acks {
				if item.IsAborted {
					session.Text = "abort"

					for _, address := range session.Participants {
						if err := m.dialer.Abort(address, session.Id); err != nil {
							log.Printf("failed to send abort message: %v\n", err)
						}
					}

					fn := time.Now()
					m.output <- session

					// update performance metrics
					du := fn.Sub(session.StartedAt).Nanoseconds() / 1000000
					m.throughput = append(m.throughput, float64(1000/du))
					m.latency = append(m.latency, float64(du))

					return
				}
			}

			// if all are committed, then commit all
			for _, address := range session.Participants {
				if err := m.dialer.Commit(address, session.Id); err != nil {
					log.Printf("failed to send commit message: %v\n", err)
				}
			}
		}
	}
}

// handleTimeouts unblocks the sessions that are not finished and they hit timeout.
func (m *Manager) handleTimeouts() {
	for {
		for key, value := range m.cache {
			if len(value.Text) == 0 && time.Since(value.StartedAt) >= 5*time.Second {
				// reset the timer
				value.StartedAt = time.Now()

				// resend the transaction
				if value.Type == "inter-shard" {
					// for inter-shard send request message to the cluster
					if err := m.dialer.Request(value.Participants[0], value.Sender, value.Receiver, value.Amount, key); err != nil {
						log.Printf("failed to resend the request for trx %d: %v\n", key, err)
					}
				} else {
					// for cross-shard send prepare messages to both clusters
					if err := m.dialer.Prepare(value.Participants[0], value.Sender, value.Sender, value.Receiver, value.Amount, key); err != nil {
						log.Printf("failed to resend the request for trx %d: %v\n", key, err)
					}
					if err := m.dialer.Prepare(value.Participants[1], value.Receiver, value.Sender, value.Receiver, value.Amount, key); err != nil {
						log.Printf("failed to resend the request for trx %d: %v\n", key, err)
					}
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
