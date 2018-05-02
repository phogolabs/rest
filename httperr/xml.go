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
	switch err.(type) {
	case xml.UnmarshalError, *xml.SyntaxError, *xml.TagPathError:
		return New(CodeInvalid, xmlUnmarshalMsg).With(http.StatusBadRequest)
	case *xml.UnsupportedTypeError:
		return New(CodeInternal, xmlMarshalMsg).With(http.StatusInternalServerError)
	default:
		return New(CodeInternal, xmlMsg).With(http.StatusInternalServerError)
	}
}
