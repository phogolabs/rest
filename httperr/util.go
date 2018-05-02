package httperr

import (
	"fmt"
	"net/http"
)

// ParamRequired is an error that occurs when the URL param is missing
func ParamRequired(key string) error {
	msg := fmt.Sprintf("Parameter '%s' is required", key)
	err := &Response{
		StatusCode: http.StatusBadRequest,
		Err:        New(CodeParamRequired, msg),
	}
	return err
}

// ParamParse is an error that occurs when the param cannot be parsed
func ParamParse(key, tname string, err error, details ...string) error {
	info := fmt.Sprintf("Parameter '%s' is not valid %s", key, tname)
	errx := &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        New(CodeParamInvalid, info, details...),
	}
	errx.Err.Wrap(err)
	return errx
}

// QueryParamRequired is an error that occurs when the URL Query param is missing
func QueryParamRequired(key string) error {
	msg := fmt.Sprintf("Query Parameter '%s' is required", key)
	err := &Response{
		StatusCode: http.StatusBadRequest,
		Err:        New(CodeQueryParamRequired, msg),
	}
	return err
}

// QueryParamParse is an error that occurs when the query param cannot be parsed
func QueryParamParse(key, tname string, err error, details ...string) error {
	info := fmt.Sprintf("Query Parameter '%s' is not valid %s", key, tname)
	errx := &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        New(CodeQueryParamInvalid, info, details...),
	}
	errx.Err.Wrap(err)
	return errx
}
