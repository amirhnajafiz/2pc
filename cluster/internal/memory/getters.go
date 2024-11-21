package memory

// GetLeader returns the current leader.
func (s *SharedMemory) GetLeader() string {
	return s.leader
}
