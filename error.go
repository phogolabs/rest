package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/form"
	"github.com/goware/errorx"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/phogolabs/rest/middleware"
	rollbar "github.com/rollbar/rollbar-go"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	// RollbarError reports errors to rollbar
	RollbarError = middleware.RollbarError

	// RollbarMessage reports errors to rollbar
	RollbarMessage = middleware.RollbarMessage
)

// Error injects the error withing the request
func Error(w http.ResponseWriter, r *http.Request, err error) {
	code := errorCode(r)

	RollbarError(rollbar.ERR, r, err)

	GetLogger(r).
		WithError(err).
		WithField("status", code).
		Error("occurred")

	Respond(w, r, WrapError(code, err))
}

// ErrorXML injects the error withing the request
func ErrorXML(w http.ResponseWriter, r *http.Request, err error) {
	code := errorCode(r)

	RollbarError(rollbar.ERR, r, err)

	GetLogger(r).
		WithError(err).
		WithField("status", code).
		Error("occurred")

	XML(w, r, WrapError(code, err))
}

// ErrorJSON injects the error withing the request
func ErrorJSON(w http.ResponseWriter, r *http.Request, err error) {
	code := errorCode(r)

	RollbarError(rollbar.ERR, r, err)

	GetLogger(r).
		WithError(err).
		WithField("status", code).
		Error("occurred")

	JSON(w, r, WrapError(code, err))
}

// ErrorStatus returns the status code for given error
func ErrorStatus(r *http.Request, err error) {
	render.Status(r, ErrorCode(err))
}

// ErrorCode returns the code for given error
func ErrorCode(err error) int {
	code := http.StatusInternalServerError

	switch err {
	case sql.ErrNoRows:
		code = http.StatusNotFound
	default:
		switch terr := err.(type) {
		case validator.ValidationErrors:
			code = http.StatusUnprocessableEntity
		case form.DecodeErrors:
			code = http.StatusBadRequest
		case *json.UnmarshalFieldError:
			code = http.StatusBadRequest
		case *json.UnmarshalTypeError:
			code = http.StatusBadRequest
		case *errorx.Errorx:
			code = terr.Code
		}
	}

	return code
}

// WrapError creates a new error
func WrapError(code int, err error) *errorx.Errorx {
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
	return errx
}

func errorCode(r *http.Request) int {
	code, ok := r.Context().Value(render.StatusCtxKey).(int)
	if !ok {
		code = http.StatusInternalServerError
		Status(r, code)
	}

	return code
}
