package memory

// GetLeader returns the current leader.
func (s *SharedMemory) GetLeader() string {
	return s.leader
}

// GetNodeName returns the current node.
func (s *SharedMemory) GetNodeName() string {
	return s.nodeName
}
