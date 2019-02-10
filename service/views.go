package service

import (
	"net/http"
	"fmt"

	"github.com/go-chi/render"
	"github.com/go-chi/chi"
)

// ShortenURL is api endpoint for creating short versions of urls
func (s *Service) ShortenURL(w http.ResponseWriter, r *http.Request) {
	data := r.Form.Get("url")
	fmt.Println(data)
}

// RedirectURL redirects from short url to its original
func (s *Service) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if shortURL := chi.URLParam(r, "shortURL"); shortURL != "" {

	} else {
		render.Render(w, r, ErrNotFound)
	}
}

// GetURLStats ...
func (s *Service) GetURLStats(w http.ResponseWriter, r *http.Request) {

}
