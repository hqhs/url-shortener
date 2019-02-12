package service

import (
	"net/http"

	// "example.com/url-shortener/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/docgen"
	kitlog "github.com/go-kit/kit/log"
	"github.com/oxtoacart/bpool"
)

// Options represents Service configuration
type Options struct {
	Domain string
	Port   string
	DbAddr string
	Driver func(string) (Database, error)
}

// Service replresent state of url-shortener service
type Service struct {
	domain string
	port   string
	r      *chi.Mux
	db     Database
	dbAddr string
	logger kitlog.Logger
	bpool  *bpool.BufferPool
}

// NewService initializes url-shortener service with database connection and url schema
func NewService(logger kitlog.Logger, o Options) (Service, error) {
	db, err := o.Driver(o.DbAddr)
	if err != nil {
		return Service{}, nil
	}
	if logger == nil {
		logger = kitlog.NewNopLogger()
	}
	s := Service{
		o.Domain,
		o.Port,
		chi.NewRouter(),
		db,
		o.DbAddr,
		logger,
		bpool.NewBufferPool(64),
	}
	s.InitRouter()
	return s, nil
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
