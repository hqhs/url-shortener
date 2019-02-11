package service

import (
	"net/http"

	// "example.com/url-shortener/redis"
	"github.com/go-chi/docgen"
	"github.com/go-chi/chi"
)

// Service replresent state of url-shortener service
type Service struct {
	r *chi.Mux
	db Database
}

// NewService initializes url-shortener service with database connection and url schema
func NewService() Service {
	// new comment
	r := chi.NewRouter()
	db := newMockDatabase()
	// db, err := redis.NewDBDriver()
	s := Service{r, db}
	s.InitRouter()
	return s
}

// Serve starts http server
func (s *Service) Serve() {
	http.ListenAndServe(":3333", s.r)
}

// RoutesDoc ...
func (s *Service) RoutesDoc() string {
	return docgen.MarkdownRoutesDoc(s.r, docgen.MarkdownOpts{
		ProjectPath: "github.com/go-chi/chi",
		Intro:       "Welcome to the chi/_examples/rest generated docs.",
	})
}
