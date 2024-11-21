package memory

// SetLeader updates the current leader id.
func (s *SharedMemory) SetLeader(leader string) {
	s.leader = leader
}
