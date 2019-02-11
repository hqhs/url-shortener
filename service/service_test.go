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

func TestApiShortening(t *testing.T) {
	service := NewService("localhost", NewMockDatabase)
	t.Run("single url", func(t *testing.T) {
		payloadURL := "https://www.google.com/search?q=golang+specification"
		payload := []byte(fmt.Sprintf(`{"url":"%s"}`, payloadURL))
		url := fmt.Sprintf("http://%s/api/v1/shorten", service.Domain)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		service.r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: %v, want: %v", status, http.StatusOK)
		}
		d := json.NewDecoder(rr.Body)
		response := URLRequest{}
		d.Decode(&response)
		// check if redirect works
		rr = httptest.NewRecorder()
		url = fmt.Sprintf("http://%s", response.RedirectURL)
		req, err = http.NewRequest(http.MethodGet, url, nil)
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

	})
}
