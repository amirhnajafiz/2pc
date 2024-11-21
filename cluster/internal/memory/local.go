package memory

// SharedMemory is a local storage for processes and handlers.
type SharedMemory struct {
	leader string
}

// NewSharedMemory returns an instance of shared memory.
func NewSharedMemory(leader string) *SharedMemory {
	return &SharedMemory{
		leader: leader,
	}
}
