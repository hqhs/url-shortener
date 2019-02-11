package service

import (
	"fmt"
	"sync"
)

// Database represent key-value store for service usage
type Database interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
}

type MockDatabase struct {
	data map[string][]byte
	mut *sync.Mutex
}


func NewMockDatabase() (Database, error) {
	data := make(map[string][]byte)
	m := &sync.Mutex{}
	return &MockDatabase{data, m}, nil
}

func (m *MockDatabase) Get(key []byte) ([]byte, error) {
	data, ok := m.data[string(key)]
	if ok {
		return data, nil
	}
	return data, fmt.Errorf("there is no value for provided key")
}

// Set tries to set value to provided key, return error if key already exists
func (m *MockDatabase) Set(key, value []byte) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	if _, ok := m.data[string(key)]; ok {
		return fmt.Errorf("provided key already has value in database")
	}
	m.data[string(key)] = value
	return nil
}
