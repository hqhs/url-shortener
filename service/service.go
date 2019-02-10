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

func NewService() Service {
	r := chi.NewRouter()
	db := newMockDatabase()
	// db := redis.NewDBDriver()
	s := Service{r, db}
	s.InitRouter()
	// Simple health check: make ping request, get or create url id
	// conn := p.Get()
	// defer conn.Close()
	// _, err := redis.String(conn.Do("PING", ""))
	// if err != nil {
	// 	panic(err)
	// }
	// id, err := GetOrCreateID(conn)
	return s
}

func (s *Service) Serve() {
	http.ListenAndServe(":3333", s.r)
}

func (s *Service) RoutesDoc() string {
	return docgen.MarkdownRoutesDoc(s.r, docgen.MarkdownOpts{
		ProjectPath: "github.com/go-chi/chi",
		Intro:       "Welcome to the chi/_examples/rest generated docs.",
	})
}
// func main() {

// 	p := NewPool("redis:6379")

// 	fmt.Printf("url id starts with: %v\n", id)

// 	r := NewRouter(p)
// 	fmt.Println("Initialized router...")


// 	http.ListenAndServe(":3333", r)
// }
