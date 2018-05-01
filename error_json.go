package rho

import (
	"encoding/json"
	"net/http"
)

const (
	jsonMsg          = "JSON Error"
	jsonUnmarshalMsg = "Unable to unmarshal json body"
	jsonMarshalMsg   = "Unable to marshal json"
)

// JSONError creates a ErrorResponse for given json error
func JSONError(err error) *ErrorResponse {
	response := &ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Err:        NewError(ErrInvalid, jsonMsg),
	}

	switch err.(type) {
	case *json.InvalidUnmarshalError:
		response.Err.Message = jsonUnmarshalMsg
	case *json.UnmarshalFieldError:
		response.Err.Message = jsonUnmarshalMsg
		response.StatusCode = http.StatusBadRequest
	case *json.UnmarshalTypeError:
		response.Err.Message = jsonUnmarshalMsg
		response.StatusCode = http.StatusBadRequest
	case *json.UnsupportedTypeError:
		response.Err.Message = jsonMarshalMsg
	case *json.UnsupportedValueError:
		response.Err.Message = jsonMarshalMsg
	case *json.InvalidUTF8Error:
		response.Err.Message = jsonMarshalMsg
	case *json.MarshalerError:
		response.Err.Message = jsonMarshalMsg
	}

	response.Err.Wrap(err)
	return response
}
