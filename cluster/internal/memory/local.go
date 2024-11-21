package memory

// SharedMemory is a local storage for processes and handlers.
type SharedMemory struct {
	leader      string
	nodeName    string
	clusterName string
	clusterIPs  []string
	iptable     map[string]string
}

// NewSharedMemory returns an instance of shared memory.
func NewSharedMemory(leader, nm, cn string, iptable map[string]string) *SharedMemory {
	instance := &SharedMemory{
		leader:      leader,
		nodeName:    nm,
		clusterName: cn,
		iptable:     iptable,
	}

	instance.SetClusterIPs()

	return instance
}
