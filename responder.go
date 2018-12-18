package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/render"
	"github.com/go-playground/form"
	"github.com/goware/errorx"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/phogolabs/rest/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	// Respond is a package-level variable set to our default Responder. We do this
	// because it allows you to set render.Respond to another function with the
	// same function signature, while also utilizing the render.Responder() function
	// itself. Effectively, allowing you to easily add your own logic to the package
	// defaults. For example, maybe you want to test if v is an error and respond
	// differently, or log something before you respond.
	Respond = DefaultResponder

	// Render renders a single payload and respond to the client request.
	Render = render.Render

	// RenderList renders a slice of payloads and responds to the client request.
	RenderList = render.RenderList

	// Status sets a HTTP response status code hint into request context at any point
	// during the request life-cycle. Before the Responder sends its response header
	// it will check the StatusCtxKey
	Status = render.Status
)

// StatusErr returns the status code for given error
func StatusErr(r *http.Request, err error) {
	code := http.StatusInternalServerError

	if err == sql.ErrNoRows {
		code = http.StatusNotFound
	}

	switch err.(type) {
	case validator.ValidationErrors:
		code = http.StatusUnprocessableEntity
	case form.DecodeErrors:
		code = http.StatusBadRequest
	case *json.UnmarshalFieldError:
		code = http.StatusBadRequest
	case *json.UnmarshalTypeError:
		code = http.StatusBadRequest
	}

	render.Status(r, code)
}

// DefaultResponder handles streaming JSON and XML responses, automatically setting the
// Content-Type based on request headers. It will default to a JSON response.
func DefaultResponder(w http.ResponseWriter, r *http.Request, v interface{}) {
	if err, ok := v.(error); ok {

		if logger := middleware.GetLogEntry(r); logger != nil {
			logger.Panic(err, debug.Stack())
		}

		code, ok := r.Context().Value(render.StatusCtxKey).(int)

		if !ok {
			code = http.StatusInternalServerError
			render.Status(r, code)
		}

		errx := errorx.New(code, http.StatusText(code))

		if errs, ok := err.(*multierror.Error); ok {
			for _, err := range errs.Errors {
				errx.Details = append(errx.Details, err.Error())
			}
		} else if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, verr := range verrs {
				if ferr, ok := verr.(error); ok {
					errx.Details = append(errx.Details, ferr.Error())
				}
			}
		} else {
			errx.Details = append(errx.Details, err.Error())
		}

		v = errx
	}

	render.DefaultResponder(w, r, v)
}
