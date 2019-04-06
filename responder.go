package rest

import (
	"net/http"

	"github.com/go-chi/render"
)

// Respond handles streaming JSON and XML responses, automatically setting the
// Content-Type based on request headers. It will default to a JSON response.
func Respond(w http.ResponseWriter, r *http.Request, v interface{}) {
	if err, ok := v.(error); ok {
		v = errorf(r, err)
	}

	// TODO: set defaults
	// TODO: set headers
	render.Respond(w, r, v)
}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	if err, ok := v.(error); ok {
		v = errorf(r, err)
	}

	// TODO: set defaults
	// TODO: set headers
	render.JSON(w, r, v)
}

// XML marshals 'v' to JSON, setting the Content-Type as application/xml. It
// will automatically prepend a generic XML header (see encoding/xml.Header) if
// one is not found in the first 100 bytes of 'v'.
func XML(w http.ResponseWriter, r *http.Request, v interface{}) {
	if err, ok := v.(error); ok {
		v = errorf(r, err)
	}

	// TODO: set defaults
	// TODO: set headers
	render.XML(w, r, v)
}

// PlainText writes a string to the response, setting the Content-Type as
// text/plain.
func PlainText(w http.ResponseWriter, r *http.Request, v string) {
	render.PlainText(w, r, v)
}

// Data writes raw bytes to the response, setting the Content-Type as
// application/octet-stream.
func Data(w http.ResponseWriter, r *http.Request, v []byte) {
	render.Data(w, r, v)
}

// HTML writes a string to the response, setting the Content-Type as text/html.
func HTML(w http.ResponseWriter, r *http.Request, v string) {
	render.HTML(w, r, v)
}

// NoContent returns a HTTP 204 "No Content" response.
func NoContent(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

// Error injects the error within the request
// Deprecated: Use Respond instead
func Error(w http.ResponseWriter, r *http.Request, err error) {
	Respond(w, r, err)
}

// ErrorXML injects the error within the request
// Deprecated: Use XML instead
func ErrorXML(w http.ResponseWriter, r *http.Request, err error) {
	XML(w, r, err)
}

// ErrorJSON injects the error within the request
// Deprecated: Use JSON instead
func ErrorJSON(w http.ResponseWriter, r *http.Request, err error) {
	JSON(w, r, err)
}
