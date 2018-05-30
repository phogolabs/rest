package httpr

import (
	"net/http"

	"github.com/go-chi/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// DefaultValidator is the default payload validator
var DefaultValidator = validator.New()

// Decode decodes a request into a struct
func Decode(r *http.Request, v interface{}) error {
	if err := render.Decode(r, v); err != nil {
		return &ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        NewError(CodeInvalid, "Unable to unmarshal request body").Wrap(err),
		}
	}

	if binder, ok := v.(render.Binder); ok {
		if err := binder.Bind(r); err != nil {
			return &ErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        NewError(CodeConditionNotMet, "Unable to bind request").Wrap(err),
			}
		}
	}

	if err := DefaultValidator.Struct(v); err != nil {
		return ValidationError(err, "Unable to validate request")
	}

	return nil
}
