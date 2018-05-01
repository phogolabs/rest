package rho

import (
	"encoding/xml"
	"net/http"
)

const (
	xmlMsg          = "XML Error"
	xmlUnmarshalMsg = "Unable to unmarshal xml body"
	xmlMarshalMsg   = "Unable to marshal xml"
)

// XMLError creates a ErrorResponse for given xml error
func XMLError(err error) *ErrorResponse {
	response := &ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Err:        NewError(ErrInternal, xmlMsg),
	}

	switch err.(type) {
	case xml.UnmarshalError:
		response.Err.Code = ErrInvalid
		response.Err.Message = xmlUnmarshalMsg
		response.StatusCode = http.StatusBadRequest
	case *xml.SyntaxError:
		response.Err.Code = ErrInvalid
		response.Err.Message = xmlUnmarshalMsg
		response.StatusCode = http.StatusBadRequest
	case *xml.TagPathError:
		response.Err.Code = ErrInvalid
		response.Err.Message = xmlUnmarshalMsg
		response.StatusCode = http.StatusBadRequest
	case *xml.UnsupportedTypeError:
		response.Err.Message = xmlMarshalMsg
	}

	response.Err.Wrap(err)
	return response
}
