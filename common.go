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

func init() {
	render.Decode = DefaultDecoder
	render.Respond = DefaultResponder
}

// RegisterValidation adds a validation with the given tag
func RegisterValidation(tag string, fn validator.Func) {
	validationFuncMap[tag] = fn
}

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

// Error injects the error withing the request
func Error(r *http.Request, err error) {
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
