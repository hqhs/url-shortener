package service

import (
	"fmt"
	"math/big"
	"sync"
)

// Database represent persistent key-value store
type Database interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	// Separate method for counter is used because
	// incrementing and getting new value should be atomic operation
	IncrementCounter() (*big.Int, error)
}

const counterKey = "counter:id"

// MockDatabase implements Database interface, not persistent, bul well-suited for tests
type MockDatabase struct {
	data       map[string][]byte
	mut        *sync.Mutex
}

// NewMockDatabase initializes MockDatabase
func NewMockDatabase(addr string) (Database, error) {
	data := make(map[string][]byte)
	m := &sync.Mutex{}
	return &MockDatabase{data, m}, nil
}

// Get takes key and return value if key is in storage, error otherwise
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

// IncrementCounter increments counter by 1 and return new value
// If there is no counter in database, return 0
func (m *MockDatabase) IncrementCounter() (*big.Int, error) {
	counter := big.NewInt(int64(0))
	m.mut.Lock()
	defer m.mut.Unlock()
	if data, ok := m.data[counterKey]; ok {
		counter = counter.SetBytes(data)
		counter = counter.Add(counter, big.NewInt(int64(0)))
		m.data[counterKey] = counter.Bytes()
		return counter, nil
	}
	m.data[counterKey] = counter.Bytes()
	return counter, nil
}
