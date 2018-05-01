package rho

import (
	"net/http"
)

// ConvError creates a ErrorResponse for given strconv error
func ConvError(err error) *ErrorResponse {
	response := &ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        NewError(ErrInvalid, "Unable to parse number"),
	}

	response.Err.Wrap(err)
	return response
}

// TimeError creates a ErrorResponse for given time error
func TimeError(err error) *ErrorResponse {
	response := &ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        NewError(ErrInvalid, "Unable to parse date time"),
	}

	response.Err.Wrap(err)
	return response
}
