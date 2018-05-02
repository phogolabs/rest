package rho

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/phogolabs/rho/httperr"
	validate "gopkg.in/go-playground/validator.v9"
)

// DefaultValidator is the default payload validator
var DefaultValidator = validate.New()

// Decode decodes a request into a struct
func Decode(r *http.Request, v interface{}) error {
	if err := render.Decode(r, v); err != nil {
		errx := httperr.New(httperr.CodeInvalid, "Unable to unmarshal request body")
		return errx.Wrap(err).With(http.StatusBadRequest)
	}

	if binder, ok := v.(render.Binder); ok {
		if err := binder.Bind(r); err != nil {
			errx := httperr.New(httperr.CodeConditionNotMet, "Unable to bind request")
			return errx.Wrap(err).With(http.StatusUnprocessableEntity)
		}
	}

	if err := DefaultValidator.Struct(v); err != nil {
		errx := httperr.New(httperr.CodeConditionNotMet, "Unable to validate request")
		return errx.Wrap(err).With(http.StatusUnprocessableEntity)
	}

	return nil
}
