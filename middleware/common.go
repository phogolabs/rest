package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
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

	// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
	// a router (or subrouter) from being cached by an upstream proxy and/or client.
	NoCache = middleware.NoCache

	// SetContentType is a middleware that forces response Content-Type.
	SetContentType = render.SetContentType

	// GetReqID returns a request ID from the given context if one is present.
	// Returns the empty string if a request ID cannot be found.
	GetReqID = middleware.GetReqID

	// Status sets a HTTP response status code hint into request context at any point
	// during the request life-cycle. Before the Responder sends its response header
	// it will check the StatusCtxKey
	Status = render.Status

	// Heartbeat endpoint middleware useful to setting up a path like
	// `/ping` that load balancers or uptime testing external services
	// can make a request before hitting any routes. It's also convenient
	// to place this above ACL middlewares as well.
	Heartbeat = middleware.Heartbeat
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
