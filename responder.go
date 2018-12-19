package rest

import (
	"github.com/go-chi/render"
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
	Respond = render.Respond

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

	// XML marshals 'v' to JSON, setting the Content-Type as application/xml. It
	// will automatically prepend a generic XML header (see encoding/xml.Header) if
	// one is not found in the first 100 bytes of 'v'.
	XML = render.XML

	// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
	// Content-Type as application/json.
	JSON = render.JSON

	// NoContent returns a HTTP 204 "No Content" response.
	NoContent = render.NoContent
)
