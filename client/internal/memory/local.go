package memory

import (
	"sync"
	"time"
)

// Memory is a local storage for client operations.
type Memory struct {
	sessions int
	lock     sync.Mutex
}

// NewMemory returns a memory instance.
func NewMemory() *Memory {
	return &Memory{
		sessions: int(time.Now().Unix()),
		lock:     sync.Mutex{},
	}
}
