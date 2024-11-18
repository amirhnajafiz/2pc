package manager

import (
	grpc "github.com/F24-CSE535/2pc/client/internal/grpc/dialer"
	"github.com/F24-CSE535/2pc/client/internal/storage"
	"github.com/F24-CSE535/2pc/client/pkg/models"
)

// Manager is a struct that handles client input commands.
type Manager struct {
	dialer  *grpc.Dialer
	storage *storage.Database

	cache map[int]*models.Session

	throughput float64
	latency    int64
	count      int
}

// NewManager returns a new manager instance.
func NewManager(dialer *grpc.Dialer, storage *storage.Database) *Manager {
	return &Manager{
		dialer:     dialer,
		storage:    storage,
		cache:      make(map[int]*models.Session),
		throughput: 0,
		latency:    0,
		count:      0,
	}
}
