package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/phogolabs/log"
	"github.com/phogolabs/rest/middleware"
)

// GetLogger returns the associated request logger
func GetLogger(r *http.Request) log.Logger {
	return middleware.GetLogger(r)
}

// Print prints the routes
func Print(routes chi.Routes) {
	chi.Walk(routes, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		fields := log.Map{
			"method": method,
			"route":  route,
		}

		log.WithFields(fields).Info("http route registered")
		return nil
	})
}
