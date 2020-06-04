package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/phogolabs/log"
)

// LoggerOption represent a logger option
type LoggerOption interface {
	Apply(logger log.Logger)
}

// LoggerOptionFunc represents a function
type LoggerOptionFunc func(logger log.Logger)

// Apply applies the option
func (fn LoggerOptionFunc) Apply(logger log.Logger) {
	fn(logger)
}

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(next http.Handler, options ...LoggerOption) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(LoggerFields(r))

		for _, option := range options {
			option.Apply(logger)
		}

		ctx := log.SetContext(r.Context(), logger)

		writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		next.ServeHTTP(writer, r.WithContext(ctx))

		logger = logger.WithFields(log.Map{
			"status":   writer.Status(),
			"size":     writer.BytesWritten(),
			"duration": time.Since(start),
		})

		switch {
		case writer.Status() >= 500:
			logger.Error("response completion fail")
		case writer.Status() >= 400:
			logger.Warn("response completion warn")
		default:
			logger.Info("response completion success")
		}
	}

	return http.HandlerFunc(fn)
}

// GetLogger returns the associated request logger
func GetLogger(r *http.Request) log.Logger {
	return log.GetContext(r.Context())
}

// LoggerFields returns the logger's fields
func LoggerFields(r *http.Request) log.Map {
	scheme := func(r *http.Request) string {
		proto := "http"

		if r.TLS != nil {
			proto = "https"
		}

		return proto
	}

	return log.Map{
		"scheme":      scheme(r),
		"host":        r.Host,
		"url":         r.RequestURI,
		"proto":       r.Proto,
		"method":      r.Method,
		"remote_addr": r.RemoteAddr,
		"request_id":  middleware.GetReqID(r.Context()),
	}
}
