package service

import (
	"net/http"
	"strings"
	"crypto/md5"
	"bytes"
	"io"
	"fmt"
	"encoding/gob"
	"encoding/base64"

	"github.com/go-chi/render"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log/level"
)

// ShortenURL is api endpoint for creating short versions of urls
func (s *Service) ShortenURL(w http.ResponseWriter, r *http.Request) {
	// add form data support?
	url := &URLRequest{}
	if err := render.Bind(r, url); err != nil {
		level.Error(s.logger).Log("url", "binding", err.Error())
		render.Render(w, r, ErrInvalidURL)
		return
	}
	buf := s.bpool.Get()
	en := gob.NewEncoder(buf)
	for {
		h := md5.New()
		counter, err := s.db.IncrementCounter()
		if err != nil {
			level.Error(s.logger).Log("increment", "counter", "error", err.Error())
			render.Render(w, r, ErrInternal)
			return
		}
		// NOTE: remove []byte -> string conversion for small perfomance boost
		io.WriteString(h, string(counter.Bytes()))
		io.WriteString(h, url.OriginalURL)
		hash := h.Sum(nil) // 16 bytes of md5 hash
		encoded := base64.RawURLEncoding.EncodeToString(hash) // url safe rfc4648 base64 encoding
		url.Key = encoded[0:6]
		if err := en.Encode(url); err != nil {
			level.Error(s.logger).Log("gob", "encoding", "url", url)
		}
		err = s.db.Set([]byte(url.Key), buf.Bytes())
		if err == nil {
			break
		}
		level.Error(s.logger).Log("key", "already", "exists", err)
		fmt.Printf("error occured: %v\n", err)
	}
	s.bpool.Put(buf)
	// FIXME are there any better solutions?
	if len(s.port) > 0 {
		url.RedirectURL = fmt.Sprintf("%s:%s/%s", s.domain, s.port, url.Key)
	} else {
		url.RedirectURL = fmt.Sprintf("%s/%s", s.domain, url.Key)
	}
	if err := render.Render(w, r, url); err != nil {
		// If service could not render it's own data, return 500 without explanation
		// Maybe add optional sentry support?
		level.Error(s.logger).Log("rendering", "short", "URL", err)
		render.Render(w, r, ErrInternal)
		return
	}
}

// RedirectURL redirects from short url to its original
func (s *Service) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if key := chi.URLParam(r, "url"); key != "" {
		encoded, err := s.db.Get([]byte(key))
		if err != nil {
			render.Render(w, r, ErrURLNotFound)
			return
		}
		decoder := gob.NewDecoder(bytes.NewReader(encoded))
		url := &URLRequest{}
		decoder.Decode(url)
		http.Redirect(w, r, url.OriginalURL, http.StatusFound)
	} else {
		// NOTE: here we could render home page
		render.Render(w, r, ErrURLNotFound)
	}
}

// GetURLStats ...
func (s *Service) GetURLStats(w http.ResponseWriter, r *http.Request) {
	payload := &Stats{}
	if err := render.Bind(r, payload); err != nil {
		level.Error(s.logger).Log("url", "binding", err.Error())
		render.Render(w, r, ErrInvalidURL)
		return
	}
	// TODO url could be localhost/HOTKEY or just HOTKEY, parse it, fetch
	// original from database, return original url with created time
	var key string
	if strings.ContainsAny(payload.URL, "/") {
		key = strings.Split(payload.URL, "/")[1]
	} else {
		key = payload.URL
	}
	encoded, err := s.db.Get([]byte(key))
	if err != nil {
		render.Render(w, r, ErrURLNotFound)
		return
	}
	url := &URLRequest{}
	decoder := gob.NewDecoder(bytes.NewReader(encoded))
	decoder.Decode(url)
	if err := render.Render(w, r, url); err != nil {
		// If service could not render it's own data, return 500 without explanation
		// Maybe add optional sentry support?
		level.Error(s.logger).Log("rendering", "short", "URL", err)
		render.Render(w, r, ErrInternal)
		return
	}
}
