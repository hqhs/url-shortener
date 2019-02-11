package service

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

//--
// Request and Response payloads for the REST api.
//--

//
type URLRequest struct {
	OriginalURL string `json:"url"`
	Key         string `json:"-"`
	RedirectURL string `json:"shortened-url"`
}

// Bind on ShortenedURL will run after the unmarshalling is complete
func (u *URLRequest) Bind(r *http.Request) error {
	if len(u.OriginalURL) == 0 {
		return fmt.Errorf("no url provided")
	}
	if _, err := url.Parse(u.OriginalURL); err != nil {
		return err
	}
	return nil
}

// Render pre-processes url before a response is marshalled and sent across the wire
func (u *URLRequest) Render(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("url: %v\n", u.RedirectURL)
	return nil
}

//--
// Error response payloads & renderers
//--

var ErrInvalidURL = &ErrResponse{HTTPStatusCode: 400, StatusText: "Provided URL is not valid."}
var ErrURLNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "URL for provided key not found."}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
