package redis

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	urlStats byte = iota
	urlIDCounter
	urlKey
)

// IdKey used for
const IdKey = "url:id"

type Redis struct {

}

func NewRedisConnection() *Redis {
	// Simple health check: make ping request, get or create url id
	// conn := p.Get()
	// defer conn.Close()
	// _, err := redis.String(conn.Do("PING", ""))
	// if err != nil {
	// 	panic(err)
	// }
	// id, err := GetOrCreateID(conn)
	return &Redis{}
}

func (r *Redis) Get(key []byte) ([]byte, error) {
	data := make([]byte, 0)
	return data, nil
}

func (r *Redis) Set(key, value []byte) (error) {
	return nil
}

// GetOrCreateID TODO
func GetOrCreateID(conn redis.Conn) (int64, error) {
	key := insertPrefixInKey(urlIDCounter, IdKey)
	id, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return redis.Int64(conn.Do("SET", key, 0))
	}
	return id, err
}

// NewPool initializes database connection
func NewPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

//--
// Data model objects and persistence mocks:
//--

// URL represents a 'Uniform Resource Locator' :)
type URL struct {
	Url   string
	Stats Stats
}

// ParseURL parses raw string into URL structure and additional data such as: ip, request time etc
func ParseURL(u string, r *http.Request) (URL, error) {
	_, err := url.Parse(u)
	if err != nil {
		return URL{}, err
	}
	stats := fetchStats(r)
	return URL{u, stats}, nil
}

// SaveURL saves given url and returns it's shorten version
func (u *URL) SaveURL(conn redis.Conn) (string, error) {
	s := u.Url
	var key []byte
	// TODO: add optional check there's no such url in database
	key = insertPrefixInKey(urlStats, u.Url)

	key = insertPrefixInKey(urlIDCounter, IdKey)
	id, _ := redis.Int64(conn.Do("GET", key))
	conn.Do("INCR", key)

	key = insertPrefixInKey(urlKey, string(id))
	_, err := redis.String(conn.Do("SET", key, s))
	if err != nil {
		return "", err
	}
	// FIXME is not url safe!
	short := encodeURL(id)
	return short, nil
}

func encodeURL(id int64) string {
	bs := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(bs, id)
	bs = bs[:n]
	return base64.StdEncoding.EncodeToString(bs)
}

// DbGetURL gets short url and decode it
func DbGetURL(short string, conn redis.Conn) (string, error) {
	id, err := decodeURL(short)
	if err != nil {
		return "", err
	}
	key := insertPrefixInKey(urlKey, string(id))
	return redis.String(conn.Do("GET", key))
}

func decodeURL(short string) (int64, error) {
	bs, err := base64.StdEncoding.DecodeString(short)
	if err != nil {
		return int64(-1), err
	}
	id, n := binary.Varint(bs)
	if n < 0 {
		return int64(-1), fmt.Errorf("value larget then 64 bits")
	} else if n == 0 {
		return int64(-1), fmt.Errorf("buf too small")
	}
	return id, nil
}

func insertPrefixInKey(prefix byte, key string) []byte {
	return append([]byte{prefix}, []byte(key)...)
}

// Stats represents shorter url usage statistics
type Stats struct {
	IP      net.IP
	Created time.Time
	Clicks  uint64
}

func fetchStats(r *http.Request) Stats {
	ip := net.ParseIP(r.RemoteAddr)
	return Stats{ip, time.Now(), uint64(0)}
}

func dbGetStats(url string) (Stats, bool) {
	return Stats{}, false
}
