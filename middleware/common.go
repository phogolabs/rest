package middleware

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var (
	// StripSlashes is a middleware that will match request paths with a trailing
	// slash, strip it from the path and continue routing through the mux, if a route
	// matches, then it will serve the handler.
	StripSlashes = middleware.StripSlashes

	// RequestID is a middleware that injects a request ID into the context of each
	// request. A request ID is a string of the form "host.example.com/random-0001",
	// where "random" is a base62 random string that uniquely identifies this go
	// process, and where the last number is an atomically incremented request
	// counter.
	RequestID = middleware.RequestID

	// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
	// of parsing either the X-Forwarded-For header or the X-Real-IP header (in that
	// order).
	RealIP = middleware.RealIP

	// DefaultCompress is a middleware that compresses response
	// body of predefined content types to a data format based
	// on Accept-Encoding request header. It uses a default
	// compression level.
	DefaultCompress = middleware.DefaultCompress

	// GetLogEntry returns the in-context LogEntry for a request.
	GetLogEntry = middleware.GetLogEntry

	// Recoverer is a middleware that recovers from panics, logs the panic (and a
	// backtrace), and returns a HTTP 500 (Internal Server Error) status if
	// possible. Recoverer prints a request ID if one is provided.
	Recoverer = middleware.Recoverer

	// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
	// a router (or subrouter) from being cached by an upstream proxy and/or client.
	NoCache = middleware.NoCache

	// SetContentType is a middleware that forces response Content-Type.
	SetContentType = render.SetContentType
)

var (
	// ErrorCtxKey is the key of the error that occurred
	ErrorCtxKey = &ContextKey{Name: "error"}
)

// ContextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type ContextKey struct {
	Name string
}

// String returns the string key
func (k *ContextKey) String() string {
	return "rest/middleware context value " + k.Name
}
