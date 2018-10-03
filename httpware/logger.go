package httpware

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
)

var (
	// DefaultLogFormatter is the default log formatter
	DefaultLogFormatter = LogFormatterFunc(NewLogEntry)
	// DefaultRequestLogger is the default logger
	DefaultRequestLogger = middleware.RequestLogger(DefaultLogFormatter)
)

var (
	_ middleware.LogFormatter = LogFormatterFunc(NewLogEntry)
	_ middleware.LogEntry     = &LogEntry{}
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(next http.Handler) http.Handler {
	return DefaultRequestLogger(next)
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

// LogEntry records the final log when a request completes.
// See defaultLogEntry for an example implementation.
type LogEntry struct {
	logger log.Interface
}

// NewLogEntry creates a new log entry
func NewLogEntry(r *http.Request) *LogEntry {
	fields := log.Fields{
		"url":        r.RequestURI,
		"proto":      r.Proto,
		"method":     r.Method,
		"remoteAddr": r.RemoteAddr,
	}

	if GetLevel() == log.DebugLevel {
		fields["header"] = r.Header
	}

	logger := log.WithFields(fields)
	return &LogEntry{logger: logger}
}

// GetLevel returns the debug level
func GetLevel() log.Level {
	if logger, ok := log.Log.(*log.Logger); ok {
		return logger.Level
	}

	if entry, ok := log.Log.(*log.Entry); ok {
		return entry.Level
	}

	return log.InvalidLevel
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
