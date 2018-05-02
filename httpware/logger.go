package httpware

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/go-chi/chi/middleware"
)

var (
	// DefaultLogFormatter is the default log formatter
	DefaultLogFormatter = LogFormatterFunc(NewLogEntry)
	// DefaultLogger is the default logger
	DefaultLogger = middleware.RequestLogger(DefaultLogFormatter)
)

var (
	_ middleware.LogFormatter = LogFormatterFunc(NewLogEntry)
	_ middleware.LogEntry     = &LogEntry{}
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}

// LogFormatterFunc is a function which implements middleware.LogFormatter
type LogFormatterFunc func(r *http.Request) *LogEntry

// NewLogEntry creates a new log entry
func (f LogFormatterFunc) NewLogEntry(r *http.Request) middleware.LogEntry {
	return f(r)
}

// LogEntry records the final log when a request completes.
// See defaultLogEntry for an example implementation.
type LogEntry struct {
	logger log.Interface
}

// NewLogEntry creates a new log entry
func NewLogEntry(r *http.Request) *LogEntry {
	logger := log.WithFields(log.Fields{
		"host":       r.Host,
		"url":        r.RequestURI,
		"proto":      r.Proto,
		"method":     r.Method,
		"remoteAddr": r.RemoteAddr,
	})

	logger.Info("request")
	return &LogEntry{logger: logger}
}

// Write logs responses
func (e *LogEntry) Write(status, bytes int, elapsed time.Duration) {
	logger := e.logger.WithFields(log.Fields{
		"status":   status,
		"size":     bytes,
		"duration": elapsed,
	})

	switch {
	case status >= 500:
		logger.Error("response")
	case status >= 400:
		logger.Warn("response")
	default:
		logger.Info("response")
	}
}

// Panic logs the panic errors
func (e *LogEntry) Panic(v interface{}, stack []byte) {
	switch err := v.(type) {
	case error:
		e.logger.WithError(err).Error("occurred")
	default:
		info := log.Fields{
			"error":  v,
			"source": string(stack),
		}
		e.logger.WithFields(info).Error("occurred")
	}
}
