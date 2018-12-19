package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/go-chi/chi/middleware"
	rollbar "github.com/rollbar/rollbar-go"
)

var (
	// DefaultLogFormatter is the default log formatter
	DefaultLogFormatter = LogFormatterFunc(NewLogEntry)

	// DefaultRequestLogger is the default logger
	DefaultRequestLogger = middleware.RequestLogger(DefaultLogFormatter)
)

var (
	_ middleware.LogEntry     = &LogEntry{}
	_ middleware.LogFormatter = LogFormatterFunc(NewLogEntry)
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(next http.Handler) http.Handler {
	return DefaultRequestLogger(next)
}

// GetLogger returns the associated request logger
func GetLogger(r *http.Request) log.Interface {
	if entry, ok := GetLogEntry(r).(*LogEntry); ok {
		return entry.logger
	}

	return log.Log
}

// SetLogger sets the logger
func SetLogger(cfg *LoggerConfig) error {
	var handler log.Handler

	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = json.New(cfg.Output)
	case "text":
		handler = text.New(cfg.Output)
	case "cli":
		handler = cli.New(cfg.Output)
	default:
		return fmt.Errorf("unsupported log format '%s'", cfg.Format)
	}

	log.SetHandler(handler)

	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}

	log.SetLevel(level)
	log.Log = log.Log.WithFields(cfg.Fields)

	return nil
}

// LogFormatterFunc is a function which implements middleware.LogFormatter
type LogFormatterFunc func(r *http.Request) *LogEntry

// NewLogEntry creates a new log entry
func (f LogFormatterFunc) NewLogEntry(r *http.Request) middleware.LogEntry {
	return f(r)
}

// LoggerConfig configures the logger
type LoggerConfig struct {
	// Fields of the root logger
	Fields log.Fields
	// Level is the logger's level (info, error, debug, verbose and etc.)
	Level string
	// Format of the log (json, text or cli)
	Format string
	// Output of the log
	Output io.Writer
}

// LogEntry records the final log when a request completes.
// See defaultLogEntry for an example implementation.
type LogEntry struct {
	request *http.Request
	logger  log.Interface
}

// NewLogEntry creates a new log entry
func NewLogEntry(request *http.Request) *LogEntry {
	return &LogEntry{
		request: request,
		logger: log.WithFields(log.Fields{
			"url":        request.RequestURI,
			"proto":      request.Proto,
			"method":     request.Method,
			"remoteAddr": request.RemoteAddr,
		}),
	}
}

// Write logs responses
func (e *LogEntry) Write(status, bytes int, elapsed time.Duration) {
	logger := e.logger.WithFields(log.Fields{
		"status":   status,
		"size":     bytes,
		"duration": elapsed,
	})

	err := fmt.Errorf(strings.ToLower(http.StatusText(status)))

	switch {
	case status >= 500:
		logger.WithError(err).Error("response")
	case status >= 400:
		logger.WithError(err).Warn("response")
	default:
		logger.Info("response")
	}
}

// Panic logs the panic errors
func (e *LogEntry) Panic(v interface{}, stack []byte) {
	info := log.Fields{
		"panic":  v,
		"source": string(stack),
	}

	e.logger.WithFields(info).Error("occurred")

	if rollbar.Token() == "" {
		return
	}

	rollbar.RequestMessageWithExtras(rollbar.CRIT, e.request, "occurred", info)
}
