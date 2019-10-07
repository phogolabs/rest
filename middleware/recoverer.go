package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/phogolabs/log"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r)

		defer func() {
			if rvr := recover(); rvr != nil {
				w.WriteHeader(http.StatusInternalServerError)

				fields := log.Map{
					"cause": rvr,
					"stack": string(debug.Stack()),
				}

				logger.WithFields(fields).Alert("panic")
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
