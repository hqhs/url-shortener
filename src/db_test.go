package main

import (
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func NewMockPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 0,
		IdleTimeout: 0,
		Dial: func() (redis.Conn, error) { return redigomock.NewConn(), nil },
	}
}

func TestMockDatabase(t *testing.T) {

}
