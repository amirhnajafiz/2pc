package manager

import (
	"fmt"
	"sort"
	"strings"

	grpc "github.com/F24-CSE535/2pc/client/internal/grpc/dialer"
	"github.com/F24-CSE535/2pc/client/internal/memory"
	"github.com/F24-CSE535/2pc/client/internal/storage"
	"github.com/F24-CSE535/2pc/client/pkg/enums"
	"github.com/F24-CSE535/2pc/client/pkg/models"
	"github.com/F24-CSE535/2pc/client/pkg/rpc/database"
)

// Manager is a struct that handles client input commands.
type Manager struct {
	dialer  *grpc.Dialer
	memory  *memory.Memory
	storage *storage.Database

	channel chan *models.Packet
	output  chan *models.Session
	cache   map[int]*models.Session

	tests map[string]*models.Testcase
	index int

	throughput []float64
	latency    []float64
}

// NewManager returns a new manager instance.
func NewManager(dialer *grpc.Dialer, storage *storage.Database) *Manager {
	// create a new manager instance
	instance := Manager{
		dialer:     dialer,
		storage:    storage,
		memory:     memory.NewMemory(),
		channel:    make(chan *models.Packet),
		output:     make(chan *models.Session),
		cache:      make(map[int]*models.Session),
		throughput: make([]float64, 0),
		latency:    make([]float64, 0),
	}

	// set default contacts
	instance.dialer.SetContacts(instance.dialer.Nodes)

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

// GetAllTests returns the current tests.
func (m *Manager) GetAllTests() ([]string, map[string]*models.Testcase) {
	// Create a slice to hold the keys
	keys := make([]string, 0, len(m.tests))

	// Populate the slice with the map's keys
	for k := range m.tests {
		keys = append(keys, k)
	}

	// Sort the slice of keys
	sort.Strings(keys)

	return keys, m.tests
}

// GetTests returns the next testcase.
func (m *Manager) GetTests() (*models.Testcase, int) {
	if m.index == len(m.tests) {
		return nil, 0
	}

	m.index++

	return m.tests[fmt.Sprintf("%d", m.index)], m.index
}

// UpdateNodesStatusForTest accepts the live-servers and contact servers and updates the manager.
func (m *Manager) UpdateNodesStatusForTest(servers []string, contacts map[string]string) error {
	// set contacts in dialer
	m.dialer.SetContacts(contacts)

	// get all servers
	all := strings.Split(m.dialer.Nodes["all"], ":")

	// block and unblock servers
	for _, s := range all {
		flag := false
		for _, ls := range servers {
			if ls == s {
				flag = true
				break
			}
		}

		if !flag {
			if err := m.dialer.Block(s); err != nil {
				return err
			}
		} else {
			if err := m.dialer.Unblock(s); err != nil {
				return err
			}
		}
	}

	return nil
}

// processor receives all gRPC messages to send the replys.
func (m *Manager) processor() {
	// start timeout handler
	go m.handleTimeouts()

	for {
		// get packets from gRPC level
		pkt := <-m.channel

		// switch case for packet label
		switch pkt.Label {
		case enums.PktReply:
			m.handleReply(pkt.Payload.(*database.ReplyMsg))
		case enums.PktAck:
			m.handleAck(pkt.Payload.(*database.AckMsg))
		}
	}
}
