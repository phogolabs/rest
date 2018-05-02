package httperr

import (
	"net/http"
)

// ConvError creates a Response for given strconv error
func ConvError(err error) *Response {
	return New(CodeInvalid, "Unable to parse number").With(http.StatusUnprocessableEntity)
}

// TimeError creates a Response for given time error
func TimeError(err error) *Response {
	return New(CodeInvalid, "Unable to parse date time").With(http.StatusUnprocessableEntity)
}
