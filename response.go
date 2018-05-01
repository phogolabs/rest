package rho

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	_ render.Renderer = &Response{}
	_ render.Renderer = &ErrorResponse{}
	_ error           = &ErrorResponse{}
)

// ResponseMeta keeps meta information for the successful response
type ResponseMeta struct {
	// Kind of the data kept
	Kind string `json:"kind,omitempty" xml:"kind,attr"`
}

// Response represents a payload of successful response
type Response struct {
	// StatusCode for the response. Default 200 OK
	StatusCode int `json:"-" xml:"-"`
	// Meta is the metadata for this response
	Meta ResponseMeta `json:"meta,omitempty" xml:"meta,omitempty"`
	// Date of this response
	Data interface{} `json:"data,omitempty" xml:"data,omitempty"`
}

// Render renders a single response and respond to the client request.
func (p *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if p.StatusCode <= 0 {
		p.StatusCode = http.StatusOK
	}

	if p.Meta.Kind == "" {
		p.Meta.Kind = typeName(p.Data)
	}

	render.Status(r, p.StatusCode)
	return nil
}

// ErrorResponse represents a HTTP error response
type ErrorResponse struct {
	StatusCode int    `json:"-" xml:"-"`
	Err        *Error `json:"error" xml:"error"`
}

// Error returns the error message from the underlying error
func (e *ErrorResponse) Error() string {
	return e.Err.Message
}

// Render renders a single error and respond to the client request.
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if e.StatusCode <= 0 {
		e.StatusCode = http.StatusInternalServerError
	}

	e.Err = e.Err.prepare()

	if e.Err.Code <= 0 {
		e.Err.Code = ErrUnknown
	}

	render.Status(r, e.StatusCode)
	return nil
}
