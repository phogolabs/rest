package httpr

import (
	"fmt"
	"net/http"

	validator "gopkg.in/go-playground/validator.v9"
)

// ValidationErrors are validation errors
type ValidationErrors = validator.ValidationErrors

// ValidationError is an error which occurrs duering validation
func ValidationError(err error, msgs ...string) *ErrorResponse {
	if len(msgs) == 0 {
		msgs = append(msgs, "Validation Error")
	}

	rerrx := NewError(CodeConditionNotMet, msgs[0], msgs[1:]...)
	response := &ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        rerrx,
	}

	errors, ok := err.(ValidationErrors)
	if !ok {
		rerrx.Reason = err
		return response
	}

	errs := MultiError{}

	for _, ferr := range errors {
		msg := fmt.Sprintf("Field '%s' is not valid", ferr.Field())
		errx := NewError(CodeFieldInvalid, msg)

		if err, ok := ferr.(error); ok {
			errx.Reason = err
		}

		errs = append(errs, errx)
	}

	rerrx.Reason = errs
	return response
}
