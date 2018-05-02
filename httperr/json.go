package httperr

import (
	"encoding/json"
	"net/http"
)

const (
	jsonMsg          = "JSON Error"
	jsonUnmarshalMsg = "Unable to unmarshal json body"
	jsonMarshalMsg   = "Unable to marshal json"
)

// JSONError creates a Response for given json error
func JSONError(err error) *Response {
	switch err.(type) {
	case *json.InvalidUnmarshalError:
		return New(CodeInternal, jsonUnmarshalMsg).With(http.StatusInternalServerError)
	case *json.UnmarshalFieldError, *json.UnmarshalTypeError:
		return New(CodeInvalid, jsonUnmarshalMsg).With(http.StatusBadRequest)
	case *json.UnsupportedTypeError, *json.UnsupportedValueError, *json.InvalidUTF8Error, *json.MarshalerError:
		return New(CodeInternal, jsonMarshalMsg).With(http.StatusInternalServerError)
	default:
		return New(CodeInternal, jsonMsg).With(http.StatusInternalServerError)
	}
}
