package manager

import (
	grpc "github.com/F24-CSE535/2pc/client/internal/grpc/dialer"
	"github.com/F24-CSE535/2pc/client/internal/storage"
	"github.com/F24-CSE535/2pc/client/pkg/models"
	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

// Manager is a struct that handles client input commands.
type Manager struct {
	dialer  *grpc.Dialer
	storage *storage.Database

	channel chan *models.Packet
	output  chan *models.Session
	cache   map[int]*models.Session

	throughput float64
	latency    int64
	count      int
}

// NewManager returns a new manager instance.
func NewManager(dialer *grpc.Dialer, storage *storage.Database) *Manager {
	// create a new manager instance
	instance := Manager{
		dialer:     dialer,
		storage:    storage,
		channel:    make(chan *models.Packet),
		output:     make(chan *models.Session),
		cache:      make(map[int]*models.Session),
		throughput: 0,
		latency:    0,
		count:      0,
	}

	// start the processor inside a go-routine
	go instance.processor()

	return &instance
}

// GetChannel returns the processor channel.
func (m *Manager) GetChannel() chan *models.Packet {
	return m.channel
}

// GetOutputChannel returns the processor ourput channel.
func (m *Manager) GetOutputChannel() chan *models.Session {
	return m.output
}

// processor receives all gRPC messages to send the replys.
func (m *Manager) processor() {
	for {
		// get packets from gRPC level
		pkt := <-m.channel

		// switch case for packet label
		switch pkt.Label {
		case "reply":
			m.handleReply(pkt.Payload.(*database.ReplyMsg))
		case "ack":
			m.handleAck(pkt.Payload.(*database.AckMsg))
		}
	}
}
