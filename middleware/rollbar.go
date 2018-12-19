package middleware

import (
	"net/http"

	"github.com/go-chi/render"
	rollbar "github.com/rollbar/rollbar-go"
)

// Rollbar is a middleware that reports errors to rollbar
func Rollbar(next http.Handler) http.Handler {
	level := func(code int) string {
		if code >= http.StatusInternalServerError {
			return rollbar.CRIT
		}

		return rollbar.ERR
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if err, ok := r.Context().Value(ErrorCtxKey).(error); ok {
			code := r.Context().Value(render.StatusCtxKey).(int)
			rollbar.RequestError(level(code), r, err)
		}
	}

	return http.HandlerFunc(fn)
}
