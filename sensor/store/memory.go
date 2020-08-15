package store

import "sync"

type MemoryStore struct {
	mux     sync.RWMutex
	message []byte
}

func (m *MemoryStore) Set(data []byte) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.message = data
	return nil
}

func (m *MemoryStore) Get() ([]byte, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.message, nil
}
