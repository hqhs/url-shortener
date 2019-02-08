package main

import (
	"time"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"encoding/binary"
	"encoding/base64"

	"github.com/gomodule/redigo/redis"
)

const (
	urlStats byte = iota
	urlIDCounter
	urlKey
)

// IdKey used for
const IdKey = "url:id"
	// append([]byte{ urlIDCounter }, []byte("urlId")...)

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
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

//--
// Data model objects and persistence mocks:
//--


// URL represents a 'Uniform Resource Locator' :)
type URL struct {
	*url.URL
}

// ParseURL parses rew string into URL structure
func ParseURL(u string) (URL, error) {
	url, err := url.Parse(u)
	if err != nil {
		return URL{}, err
	}
	return URL{url}, nil
}

// SaveURL saves given url and returns it's shorten version
func (u *URL) SaveURL(conn redis.Conn) (string, error) {
	// TODO: check there's no such url in database
	s := u.String()
	var key []byte

	key = insertPrefixInKey(urlIDCounter, IdKey)
	id, _ := redis.Int64(conn.Do("GET", key))
	conn.Do("INCR", key)

	key = insertPrefixInKey(urlKey, string(id))
	_, err := redis.String(conn.Do("SET", key, s))
	if err != nil {
		return "", err
	}
	short := encodeURL(id)
	return short, nil
}

func encodeURL(id int64) string {
	// So, encoding algorithm requires managing of a LOT of corner
	// cases, such as: spam filtering, no nasty words in shorten url versions, etc
	// I skipped all of them and just do base64 encoding of id
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
	created time.Time
}

// EXAMPLE CODE BELOW:

// User data model
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Article data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Article struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"` // the author
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

// Article fixture data
var articles = []*Article{
	{ID: "1", UserID: 100, Title: "Hi", Slug: "hi"},
	{ID: "2", UserID: 200, Title: "sup", Slug: "sup"},
	{ID: "3", UserID: 300, Title: "alo", Slug: "alo"},
	{ID: "4", UserID: 400, Title: "bonjour", Slug: "bonjour"},
	{ID: "5", UserID: 500, Title: "whats up", Slug: "whats-up"},
}

// User fixture data
var users = []*User{
	{ID: 100, Name: "Peter"},
	{ID: 200, Name: "Julia"},
}

func dbNewArticle(article *Article) (string, error) {
	article.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	articles = append(articles, article)
	return article.ID, nil
}

func dbGetArticle(id string) (*Article, error) {
	for _, a := range articles {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("article not found.")
}

func dbGetArticleBySlug(slug string) (*Article, error) {
	for _, a := range articles {
		if a.Slug == slug {
			return a, nil
		}
	}
	return nil, errors.New("article not found.")
}

func dbUpdateArticle(id string, article *Article) (*Article, error) {
	for i, a := range articles {
		if a.ID == id {
			articles[i] = article
			return article, nil
		}
	}
	return nil, errors.New("article not found.")
}

func dbRemoveArticle(id string) (*Article, error) {
	for i, a := range articles {
		if a.ID == id {
			articles = append((articles)[:i], (articles)[i+1:]...)
			return a, nil
		}
	}
	return nil, errors.New("article not found.")
}

func dbGetUser(id int64) (*User, error) {
	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found.")
}
