package httperr

import (
	"fmt"
	"strings"

	"github.com/gosuri/uitable"
)

const (
	// CodeParamRequired is an error code returned when the parameter is missing
	CodeParamRequired = 20101
	// CodeParamInvalid is an error code returned when the parameter's value is an invalid
	CodeParamInvalid = 20102
	// CodeQueryParamRequired is an error code returned when the query parameter is missing
	CodeQueryParamRequired = 20101
	// CodeQueryParamInvalid is an error code returned when the query parameter's value is an invalid
	CodeQueryParamInvalid = 20102
)

const (
	// CodeConflict when the API request cannot be completed because the requested operation would conflict with an existing item.
	CodeConflict = 40101
	// CodeDuplicate when the requested operation failed because it tried to create a resource that already exists.
	CodeDuplicate = 40102
	// CodeDeleted when the request failed because the resource associated with the request has been deleted 410
	CodeDeleted = 40103
	// CodeConditionNotMet when the condition set in the request's was not met 416
	CodeConditionNotMet = 40104
	// CodeOutOfrange when the request specified a range that cannot be satisfied 428.
	CodeOutOfrange = 40105
)

const (
	// CodeInternal when the request failed due to an internal error 500.
	CodeInternal = 40106
	// CodeInvalid when the provided payload is invalid
	CodeInvalid = 40222
	// CodeFieldInvalid when the struct field is invalid
	CodeFieldInvalid = 40222
	// CodeBackend when a backend error occurred 503. Usually database.
	CodeBackend = 40106
	// CodeBackendNotConnected when the request failed due to a connection error.
	CodeBackendNotConnected = 40107
	// CodeBackendNotReady code when the API server is not ready to accept requests.
	CodeBackendNotReady = 40108
)

var _ error = &Error{}

// Error is a more feature rich implementation of error interface inspired
// by PostgreSQL error style guide
type Error struct {
	Code    int      `json:"code,omitempty" xml:"code,omitempty"`
	Message string   `json:"message" xml:"message"`
	Reason  error    `json:"reason,omitempty" xml:"reason,omitempty"`
	Details []string `json:"details,omitempty" xml:"details,omitempty"`
}

// New returns an error with error code and error messages provided in
// function params
func New(code int, msg ...string) *Error {
	e := Error{Code: code}

	count := len(msg)
	if count > 0 {
		e.Message = msg[0]
	}

	if count > 1 {
		e.Details = msg[1:]
	}

	return &e
}

// With returns the error as Response Error
func (e *Error) With(status int) *Response {
	return &Response{
		StatusCode: status,
		Err:        e,
	}
}

// Error returns the error message
func (e *Error) Error() string {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	table.AddRow("code:", fmt.Sprintf("%d", e.Code))
	table.AddRow("message:", e.Message)

	if len(e.Details) > 0 {
		table.AddRow("details:", strings.Join(e.Details, ", "))
	}

	if e.Reason != nil {
		table.AddRow("reason:", e.Reason.Error())
	}

	return table.String()
}

// Wrap wraps the actual error
func (e *Error) Wrap(err error) *Error {
	e.Reason = err
	return e
}

// Cause returns the real reason for the error
func (e *Error) Cause() error {
	if e.Reason == nil {
		return e
	}

	if reason, ok := e.Reason.(*Error); ok {
		return reason.Cause()
	}

	return e.Reason
}

func (e Error) prepare() *Error {
	err := &Error{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}

	if e.Reason == nil {
		return err
	}

	if reason, ok := e.Reason.(*Error); ok {
		err.Reason = reason.prepare()
	} else {
		err.Reason = &Error{Message: e.Reason.Error()}
	}

	return err
}
