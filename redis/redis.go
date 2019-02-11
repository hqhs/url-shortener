package redis

import (
	"math/big"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	urlStats byte = iota
	urlIDCounter
	urlKey
)

const counterKey = "counter:id"

// Redis implements service.Database interface for redis connection pool
// NOTE currently every call to /api/v1/shorten uses minimum two database
// calls without connection caching. Fix it for perfomance imrpovement
type Redis struct {
	pool       *redis.Pool
	counterKey []byte
}

// NewRedisConnectionPool initializes reddis connection pool and do simple health check
func NewRedisConnectionPool(addr string) (*Redis, error) {
	// Simple health check: make ping request, get or create url id
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
	conn := pool.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("PING", ""))
	if err != nil {
		return &Redis{}, err
	}
	return &Redis{pool, []byte(counterKey)}, nil
}

// Get takes key and return value if key is in storage, error otherwise
func (r *Redis) Get(key []byte) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("GET", key))
}

// Set tries to set value to provided key, return error if key already exists
func (r *Redis) Set(key, value []byte) error {
	conn := r.pool.Get()
	defer conn.Close()
	// https://redis.io/commands/setnx
	i, err := redis.Int(conn.Do("SETNX", key, value))
	if err != nil {
		return err
	}
	if i == 0 {
		return fmt.Errorf("Key already exists")
	}
	return nil
}

// IncrementCounter increments counter by 1 and return new value
// If there is no counter in database, return 0
func (r *Redis) IncrementCounter() (*big.Int, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return big.NewInt(int64(0)), nil
}
