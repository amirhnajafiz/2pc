package lock

import "sync"

// Manager is a centeral lock manager to handle multiple accesses on objects.
type Manager struct {
	locks map[string]int
	lock  sync.Mutex
}

// NewManager returns an instance of lock manager.
func NewManager() *Manager {
	return &Manager{
		locks: make(map[string]int),
		lock:  sync.Mutex{},
	}
}
