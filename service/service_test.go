package service

import (
	"fmt"
	"testing"
	"strings"
	"bytes"
	"net/http"
	"net/http/httptest"
	"encoding/json"
)

func short(service *Service, u string) (*httptest.ResponseRecorder, error) {
	payload := []byte(fmt.Sprintf(`{"url":"%s"}`, u))
	url := fmt.Sprintf("http://%s:%s/api/v1/shorten", service.domain, service.port)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return rr, err
	}
	req.Header.Set("Content-Type", "application/json")
	service.r.ServeHTTP(rr, req)
	return rr, err
}

func TestApiShortening(t *testing.T) {
	options := Options{"localhost", "", "", NewMockDatabase}
	service, _ := NewService(nil, options)
	t.Run("single url", func(t *testing.T) {
		payloadURL := "https://www.google.com/search?q=golang+specification"
		rr, err := short(&service, payloadURL)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: %v, want: %v", status, http.StatusOK)
		}
		response := URLRequest{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		// check if redirect works
		rr = httptest.NewRecorder()
		url := fmt.Sprintf("http://%s", response.RedirectURL)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		service.r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusFound {
			t.Errorf("handler returned wrong status code: %v, want: %v", status, http.StatusFound)
		}
		r := rr.Result()
		if loc, err := r.Location(); err != nil || strings.Compare(loc.String(), payloadURL) != 0 {
			t.Errorf("handler redirected to wrong url")
		}
	})
	t.Run("invalid url", func(t *testing.T) {
		payloadURL := "notAvalidURL"
		rr, err := short(&service, payloadURL)
		if err != nil {
			t.Fatal(err)
		}
		if status := rr.Code; status == http.StatusOK {
			t.Errorf("handler shortened invalid url: %v", payloadURL)
		}
	})
}

func BenchmarkAPIShorten(b *testing.B) {
	options := Options{"localhost", "", "", NewMockDatabase}
	service, nil := NewService(nil, options)
	url := "http://www.reddit.com"
	for n := 0; n < b.N; n++ {
		if _, err := short(&service, url); err != nil {
			// TODO: disable chi logging
			b.Errorf("Unexpected error during benchmarking: %v", err)
		}
	}
}
