package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/form"
	"github.com/goware/errorx"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/phogolabs/rest/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	// Status sets a HTTP response status code hint into request context at any point
	// during the request life-cycle. Before the Responder sends its response header
	// it will check the StatusCtxKey
	Status = render.Status

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

	// PlainText writes a string to the response, setting the Content-Type as
	// text/plain.
	PlainText = render.PlainText

	// Data writes raw bytes to the response, setting the Content-Type as
	// application/octet-stream.
	Data = render.Data

	// HTML writes a string to the response, setting the Content-Type as text/html.
	HTML = render.HTML

	// NoContent returns a HTTP 204 "No Content" response.
	NoContent = render.NoContent
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

	*r = *r.WithContext(context.WithValue(r.Context(), middleware.ErrorCtxKey, err))
}

// WrapError creates a new error
func WrapError(err error) *errorx.Errorx {
	errx := errorx.New(0, "")

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
	return errx
}

// DefaultResponder handles streaming JSON and XML responses, automatically setting the
// Content-Type based on request headers. It will default to a JSON response.
func DefaultResponder(w http.ResponseWriter, r *http.Request, v interface{}) {
	v = response(r, v)
	render.DefaultResponder(w, r, v)
}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	v = response(r, v)
	render.JSON(w, r, v)
}

// XML marshals 'v' to JSON, setting the Content-Type as application/xml. It
// will automatically prepend a generic XML header (see encoding/xml.Header) if
// one is not found in the first 100 bytes of 'v'.
func XML(w http.ResponseWriter, r *http.Request, v interface{}) {
	v = response(r, v)
	render.XML(w, r, v)
}

func response(r *http.Request, v interface{}) interface{} {
	if err, ok := v.(error); ok {
		errx, ok := err.(*errorx.Errorx)

		if !ok {
			errx = WrapError(err)
		}

		if errx.Code == 0 {
			code, ok := r.Context().Value(render.StatusCtxKey).(int)

			if !ok {
				code = http.StatusInternalServerError
			}

			errx.Code = code
		}

		if errx.Message == "" {
			errx.Message = http.StatusText(errx.Code)
		}

		render.Status(r, errx.Code)
		return errx
	}

	return v
}
