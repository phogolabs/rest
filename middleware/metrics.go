package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics enables metrics for each request
func Metrics(next http.Handler) http.Handler {
	labels := []string{
		"request_id",
		"method",
		"url",
		"code",
	}

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "rest",
		Subsystem: "http",
		Name:      "request_count",
		Help:      "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
	}, labels[0:3])

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "rest",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, labels)

	fn := func(w http.ResponseWriter, r *http.Request) {
		id := middleware.GetReqID(r.Context())
		writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		gauge := gauge.WithLabelValues(id, r.Method, r.RequestURI)
		gauge.Add(1)

		start := time.Now()

		next.ServeHTTP(writer, r)

		duration := time.Since(start).Seconds()
		code := http.StatusText(writer.Status())

		histogram.WithLabelValues(id, r.Method, r.RequestURI, code).Observe(duration)

		gauge.Dec()
	}

	return http.HandlerFunc(fn)
}
