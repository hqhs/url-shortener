package service

import (
	"net/http"
	"crypto/md5"
	"io"
	"path"
	"encoding/base64"

	"github.com/go-chi/render"
	"github.com/go-chi/chi"
)

// ShortenURL is api endpoint for creating short versions of urls
func (s *Service) ShortenURL(w http.ResponseWriter, r *http.Request) {
	// add form data support?
	url := &URLRequest{}
	if err := render.Bind(r, url); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	counter, err := s.db.IncrementCounter()
	for {
		h := md5.New()
		if err != nil {
			// TODO: return 500, log error
			render.Render(w, r, ErrInternal)
			return
		}
		// NOTE: remove []byte -> string conversion for small perfomance boost
		io.WriteString(h, string(counter.Bytes()))
		io.WriteString(h, url.OriginalURL)
		hash := h.Sum(nil) // 16 bytes of md5 hash
		encoded := base64.RawURLEncoding.EncodeToString(hash) // url safe rfc4648 base64 encoding
		url.Key = encoded[0:6]
		err = s.db.Set([]byte(url.Key), []byte(url.OriginalURL))
		if err == nil {
			break
		}
	}
	url.RedirectURL = path.Join(s.Domain, url.Key) // FIXME not a good idea
	if err := render.Render(w, r, url); err != nil {
		// If service could not render it's own data, return 500 without explanation for client
		// TODO: log error, at least. Or add optional sentry support
		render.Render(w, r, ErrInternal)
		return
	}
}

// RedirectURL redirects from short url to its original
func (s *Service) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if shortURL := chi.URLParam(r, "url"); shortURL != "" {
		url, err := s.db.Get([]byte(shortURL))
		if err != nil {
			render.Render(w, r, ErrURLNotFound)
			return
		}
		http.Redirect(w, r, string(url), http.StatusFound)
	} else {
		// NOTE: here we could render home page
		render.Render(w, r, ErrURLNotFound)
	}
}

// GetURLStats ...
func (s *Service) GetURLStats(w http.ResponseWriter, r *http.Request) {

}
