package httperr

import (
	"net/http"
)

// ConvError creates a Response for given strconv error
func ConvError(err error) *Response {
	response := &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        New(ErrInvalid, "Unable to parse number"),
	}

	response.Err.Wrap(err)
	return response
}

// TimeError creates a Response for given time error
func TimeError(err error) *Response {
	response := &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        New(ErrInvalid, "Unable to parse date time"),
	}

	response.Err.Wrap(err)
	return response
}
