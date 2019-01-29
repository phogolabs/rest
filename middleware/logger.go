package middleware

import (
	"context"
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

// LoggerCtxKey is the context.Context key to store the request log entry.
var LoggerCtxKey = &ContextKey{"Logger"}

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

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(next http.Handler) http.Handler {
	scheme := func(r *http.Request) string {
		proto := "http"

		if r.TLS != nil {
			proto = "https"
		}

		return proto
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(log.Fields{
			"scheme":      scheme(r),
			"host":        r.Host,
			"url":         r.RequestURI,
			"proto":       r.Proto,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"request_id":  middleware.GetReqID(r.Context()),
		})

		writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ctx := context.WithValue(r.Context(), LoggerCtxKey, logger)

		start := time.Now()
		next.ServeHTTP(writer, r.WithContext(ctx))

		logger = logger.WithFields(log.Fields{
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
func GetLogger(r *http.Request) log.Interface {
	if logger, ok := r.Context().Value(LoggerCtxKey).(log.Interface); ok {
		return logger
	}

	return log.Log
}
