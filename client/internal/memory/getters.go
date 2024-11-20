package memory

// GetSession returns a unique session id.
func (m *Memory) GetSession() int {
	m.lock.Lock()
	session := m.sessions
	m.sessions++
	m.lock.Unlock()

	return session
}
