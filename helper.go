package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/phogolabs/log"
	"github.com/phogolabs/rest/middleware"
)

// GetLogger returns the associated request logger
func GetLogger(r *http.Request) log.Writer {
	return middleware.GetLogger(r)
}

// Print print the routes
func Print(routes chi.Routes) {
	chi.Walk(routes, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		fields := log.FieldMap{
			"method": method,
			"route":  route,
		}

		log.WithFields(fields).Info("http route registered")
		return nil
	})
}
