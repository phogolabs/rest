package httperr

import (
	"fmt"
	"net/http"

	validator "gopkg.in/go-playground/validator.v9"
)

// ValidationErrors are validation errors
type ValidationErrors = validator.ValidationErrors

// ValidationError is an error which occurrs duering validation
func ValidationError(err error, msgs ...string) *Response {
	if len(msgs) == 0 {
		msgs = append(msgs, "Validation Error")
	}

	rerrx := New(CodeConditionNotMet, msgs[0], msgs[1:]...)
	errors, ok := err.(ValidationErrors)
	if !ok {
		rerrx.Reason = err
		return rerrx.With(http.StatusUnprocessableEntity)
	}

	errs := MultiError{}

	for _, ferr := range errors {
		msg := fmt.Sprintf("Field '%s' is not valid", ferr.Field())
		errx := New(CodeFieldInvalid, msg)

		if err, ok := ferr.(error); ok {
			errx.Reason = err
		}

		errs = append(errs, errx)
	}

	rerrx.Reason = errs

	return rerrx.With(http.StatusUnprocessableEntity)
}
