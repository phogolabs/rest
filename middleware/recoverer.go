package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/apex/log"
	rollbar "github.com/rollbar/rollbar-go"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				w.WriteHeader(http.StatusInternalServerError)

				fields := log.Fields{
					"panic": rvr,
					"stack": debug.Stack(),
				}

				GetLogger(r).WithFields(fields).Error("panic")

				if rollbar.Token() == "" {
					return
				}

				rollbar.RequestMessageWithExtras(rollbar.CRIT, r, "occurred", fields)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
