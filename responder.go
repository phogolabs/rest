package rest

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/goware/errorx"
)

var (
	// Status sets a HTTP response status code hint into request context at any point
	// during the request life-cycle. Before the Responder sends its response header
	// it will check the StatusCtxKey
	Status = render.Status

	// Respond is a package-level variable set to our default Responder. We do this
	// because it allows you to set render.Respond to another function with the
	// same function signature, while also utilizing the render.Responder() function
	// itself. Effectively, allowing you to easily add your own logic to the package
	// defaults. For example, maybe you want to test if v is an error and respond
	// differently, or log something before you respond.
	Respond = DefaultResponder

	// Render renders a single payload and respond to the client request.
	Render = render.Render

	// RenderList renders a slice of payloads and responds to the client request.
	RenderList = render.RenderList

	// PlainText writes a string to the response, setting the Content-Type as
	// text/plain.
	PlainText = render.PlainText

	// Data writes raw bytes to the response, setting the Content-Type as
	// application/octet-stream.
	Data = render.Data

	// HTML writes a string to the response, setting the Content-Type as text/html.
	HTML = render.HTML

	// NoContent returns a HTTP 204 "No Content" response.
	NoContent = render.NoContent
)

// DefaultResponder handles streaming JSON and XML responses, automatically setting the
// Content-Type based on request headers. It will default to a JSON response.
func DefaultResponder(w http.ResponseWriter, r *http.Request, v interface{}) {
	v = response(r, v)
	render.DefaultResponder(w, r, v)
}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	v = response(r, v)
	render.JSON(w, r, v)
}

// XML marshals 'v' to JSON, setting the Content-Type as application/xml. It
// will automatically prepend a generic XML header (see encoding/xml.Header) if
// one is not found in the first 100 bytes of 'v'.
func XML(w http.ResponseWriter, r *http.Request, v interface{}) {
	v = response(r, v)
	render.XML(w, r, v)
}

func response(r *http.Request, v interface{}) interface{} {
	if err, ok := v.(error); ok {
		errx, ok := err.(*errorx.Errorx)

		if !ok {
			errx = WrapError(err)
		}

		if errx.Code == 0 {
			code, ok := r.Context().Value(render.StatusCtxKey).(int)

			if !ok {
				code = http.StatusInternalServerError
			}

			errx.Code = code
		}

		if errx.Message == "" {
			errx.Message = http.StatusText(errx.Code)
		}

		render.Status(r, errx.Code)
		return errx
	}

	return v
}
