package httperr

import (
	"encoding/xml"
	"net/http"
)

const (
	xmlMsg          = "XML Error"
	xmlUnmarshalMsg = "Unable to unmarshal xml body"
	xmlMarshalMsg   = "Unable to marshal xml"
)

// XMLError creates a Response for given xml error
func XMLError(err error) *Response {
	response := &Response{
		StatusCode: http.StatusInternalServerError,
		Err:        New(ErrInternal, xmlMsg),
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
