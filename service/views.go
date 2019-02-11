package service

import (
	"net/http"
	"net/url"
	"crypto/md5"
	"strconv"
	"io"
	"time"
	"fmt"
	"encoding/base64"

	"github.com/go-chi/render"
	"github.com/go-chi/chi"
)

// ShortenURL is api endpoint for creating short versions of urls
func (s *Service) ShortenURL(w http.ResponseWriter, r *http.Request) {
	data := r.Form.Get("url")
	url, error := url.Parse(data)
	if error != nil {
		render.Render(w, r, ErrInvalidURL)
		return
	}
	h := md5.New()
	now := time.Now().UnixNano()
	// FIXME there's still some probability what we got duplicate shortens
	// for same link, use counter instead
	io.WriteString(h, strconv.FormatInt(now, 10))
	io.WriteString(h, url.String())
	hash := h.Sum(nil) // 16 bytes of md5 hash
	encoded := base64.StdEncoding.EncodeToString(hash)
	short := encoded[0:6]
	// TODO repeat until short url is unique

	fmt.Println(short)
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
