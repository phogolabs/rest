package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics enables metrics for each request
func Metrics(next http.Handler) http.Handler {
	labels := []string{
		"code",
		"method",
	}

	connTotal := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_http_client_connected",
		Help: "Number of active client connections",
	})

	reqTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "go_http_requests_total",
		Help: "total HTTP requests processed",
	}, labels)

	respTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "go_http_response_duration_seconds",
		Help:    "Histogram of response time for handler",
		Buckets: prometheus.DefBuckets,
	}, labels)

	headerTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "go_http_response_header_duration_seconds",
		Help:    "Histogram of header write time for handler",
		Buckets: prometheus.DefBuckets,
	}, labels)

	return promhttp.InstrumentHandlerInFlight(connTotal,
		promhttp.InstrumentHandlerCounter(reqTotal,
			promhttp.InstrumentHandlerDuration(respTime,
				promhttp.InstrumentHandlerTimeToWriteHeader(headerTime,
					next))))
}
