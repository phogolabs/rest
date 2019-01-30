package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics enables metrics for each request
func Metrics(next http.Handler) http.Handler {
	labels := []string{
		"id",
		"code",
		"handler",
		"method",
	}

	reqInflight := promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "The HTTP requests in flight",
	})

	reqTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests made",
	}, labels[1:])

	reqTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "The HTTP response duration time",
		Buckets:   prometheus.DefBuckets,
	}, labels)

	hn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(writer, r)
		Status(r, writer.Status())
	})

	fn := func(w http.ResponseWriter, r *http.Request) {
		handler := InstrumentHandlerInFlight(reqInflight,
			InstrumentHandlerCounter(reqTotal,
				InstrumentHandlerDuration(reqTime, hn)))

		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// InstrumentHandlerInFlight is a middleware that wraps the provided
// http.Handler. It sets the provided prometheus.Gauge to the number of
// requests currently handled by the wrapped http.Handler.
//
// See the example for InstrumentHandlerDuration for example usage.
func InstrumentHandlerInFlight(gauge prometheus.Gauge, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gauge.Inc()
		defer gauge.Dec()
		next.ServeHTTP(w, r)
	})
}

// InstrumentHandlerCounter is a middleware that wraps the provided http.Handler
// to observe the request result with the provided CounterVec.  The CounterVec
// must have zero, one, or two non-const non-curried labels. For those, the only
// allowed label names are "code" and "method". The function panics
// otherwise. Partitioning of the CounterVec happens by HTTP status code and/or
// HTTP method if the respective instance label names are present in the
// CounterVec. For unpartitioned counting, use a CounterVec with zero labels.
//
// If the wrapped Handler does not set a status code, a status code of 200 is assumed.
//
// If the wrapped Handler panics, the Counter is not incremented.
//
// See the example for InstrumentHandlerDuration for example usage.
func InstrumentHandlerCounter(counter *prometheus.CounterVec, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		counter.With(InstrumentLabels(r, "code")).Inc()
	})
}

// InstrumentHandlerDuration is a middleware that wraps the provided
// http.Handler to observe the request duration with the provided ObserverVec.
// The ObserverVec must have zero, one, or two non-const non-curried labels. For
// those, the only allowed label names are "code" and "method". The function
// panics otherwise. The Observe method of the Observer in the ObserverVec is
// called with the request duration in seconds. Partitioning happens by HTTP
// status code and/or HTTP method if the respective instance label names are
// present in the ObserverVec. For unpartitioned observations, use an
// ObserverVec with zero labels. Note that partitioning of Histograms is
// expensive and should be used judiciously.
//
// If the wrapped Handler does not set a status code, a status code of 200 is assumed.
//
// If the wrapped Handler panics, no values are reported.
//
// Note that this method is only guaranteed to never observe negative durations
// if used with Go1.9+.
func InstrumentHandlerDuration(obs prometheus.ObserverVec, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		obs.With(InstrumentLabels(r, "id", "code")).Observe(time.Since(now).Seconds())
	})
}

// InstrumentLabels returns the instrument labels
func InstrumentLabels(r *http.Request, keys ...string) prometheus.Labels {
	ctx := r.Context()
	rctx := chi.RouteContext(ctx)

	labels := prometheus.Labels{
		"handler": rctx.RoutePattern(),
		"method":  rctx.RouteMethod,
	}

	for _, key := range keys {
		value := ""

		switch key {
		case "id":
			value = middleware.GetReqID(ctx)
		case "code":
			status, ok := r.Context().Value(render.StatusCtxKey).(int)
			if !ok {
				status = 0
			}
			value = fmt.Sprintf("%v", status)
		}

		labels[key] = value
	}

	return labels
}
