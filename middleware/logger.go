package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/phogolabs/log"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(LoggerFields(r))
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
			logger.Error("response")
		case writer.Status() >= 400:
			logger.Warn("response")
		default:
			logger.Info("response")
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
