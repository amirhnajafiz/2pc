package lock

// Lock returns true if a record is available and can be used.
func (m *Manager) Lock(key string, sessionId int) bool {
	m.lock.Lock()
	result := true

	// check if lock exists or not
	if _, ok := m.locks[key]; ok {
		result = false
	} else {
		result = true
		m.locks[key] = sessionId
	}
	m.lock.Unlock()

	return result
}

// Unlock releases a record.
func (m *Manager) Unlock(key string, sessionId int) {
	m.lock.Lock()
	if value, ok := m.locks[key]; ok && value == sessionId {
		delete(m.locks, key)
	}
	m.lock.Unlock()
}
