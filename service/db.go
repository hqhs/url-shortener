package service

import (
	"errors"
	"fmt"
	"math/rand"
)

// Database
type Database interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
}

type mockDatabase struct {
	data map[string][]byte
}

func newMockDatabase() *mockDatabase {
	data := make(map[string][]byte)
	return &mockDatabase{data}
}

func (m *mockDatabase) Get(key []byte) ([]byte, error) {
	data, _ := m.data[string(key)]
	return data, nil
}

func (m *mockDatabase) Set(key, value []byte) error {
	m.data[string(key)] = value
	return nil
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