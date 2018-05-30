package httpr

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// ConvError creates a Response for given strconv error
func ConvError(err error) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        NewError(CodeInvalid, "Unable to parse number").Wrap(err),
	}
}

// TimeError creates a Response for given time error
func TimeError(err error) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        NewError(CodeInvalid, "Unable to parse date time").Wrap(err),
	}
}

// JSONError creates a Response for given json error
func JSONError(err error) *ErrorResponse {
	const (
		jsonMsg          = "JSON Error"
		jsonUnmarshalMsg = "Unable to unmarshal json body"
		jsonMarshalMsg   = "Unable to marshal json"
	)

	switch err.(type) {
	case *json.InvalidUnmarshalError:
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeInternal, jsonUnmarshalMsg).Wrap(err),
		}
	case *json.UnmarshalFieldError, *json.UnmarshalTypeError:
		return &ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        NewError(CodeInvalid, jsonUnmarshalMsg).Wrap(err),
		}
	case *json.UnsupportedTypeError, *json.UnsupportedValueError, *json.InvalidUTF8Error, *json.MarshalerError:
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeInternal, jsonMarshalMsg).Wrap(err),
		}
	default:
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeInternal, jsonMsg).Wrap(err),
		}
	}
}

// XMLError creates a Response for given xml error
func XMLError(err error) *ErrorResponse {
	const (
		xmlMsg          = "XML Error"
		xmlUnmarshalMsg = "Unable to unmarshal xml body"
		xmlMarshalMsg   = "Unable to marshal xml"
	)

	switch err.(type) {
	case xml.UnmarshalError, *xml.SyntaxError, *xml.TagPathError:
		return &ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        NewError(CodeInvalid, xmlUnmarshalMsg).Wrap(err),
		}
	case *xml.UnsupportedTypeError:
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeInternal, xmlMarshalMsg).Wrap(err),
		}
	default:
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeInternal, xmlMsg).Wrap(err),
		}
	}
}

// UnknownError handles httperr
func UnknownError(err error) *ErrorResponse {
	if err == nil {
		return nil
	}

	switch errx := err.(type) {
	case *ErrorResponse:
		return errx
	case *HTTPError:
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errx,
		}
	default:
		return &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeInternal, "Internal Error").Wrap(err),
		}
	}
}
