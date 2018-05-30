package httpr

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

const code = 10000

const (
	// CodeParamRequired is an error code returned when the parameter is missing
	CodeParamRequired = iota + code
	// CodeParamInvalid is an error code returned when the parameter's value is an invalid
	CodeParamInvalid = iota + code
	// CodeQueryParamRequired is an error code returned when the query parameter is missing
	CodeQueryParamRequired = iota + code
	// CodeQueryParamInvalid is an error code returned when the query parameter's value is an invalid
	CodeQueryParamInvalid = iota + code
	// CodeConflict when the API request cannot be completed because the requested operation would conflict with an existing item.
	CodeConflict = iota + code
	// CodeDuplicate when the requested operation failed because it tried to create a resource that already exists.
	CodeDuplicate = iota + code
	// CodeDeleted when the request failed because the resource associated with the request has been deleted 410
	CodeDeleted = iota + code
	// CodeConditionNotMet when the condition set in the request's was not met 416
	CodeConditionNotMet = iota + code
	// CodeOutOfrange when the request specified a range that cannot be satisfied 428.
	CodeOutOfrange = iota + code
	// CodeInternal when the request failed due to an internal error 500.
	CodeInternal = iota + code
	// CodeInvalid when the provided payload is invalid
	CodeInvalid = iota + code
	// CodeFieldInvalid when the struct field is invalid
	CodeFieldInvalid = iota + code
	// CodeBackend when a backend error occurred 503. Usually database.
	CodeBackend = iota + code
	// CodeBackendNotConnected when the request failed due to a connection error.
	CodeBackendNotConnected = iota + code
	// CodeBackendNotReady code when the API server is not ready to accept requests.
	CodeBackendNotReady = iota + code
)

// Error represents an error that can be prepared for rendering
type Error interface {
	// Error returns the error message
	Error() string
	// Prepare prepares the error
	Prepare() Error
	// Fields returns the fields used by the logger
	Fields() log.Fields
}

var (
	_ Error = &HTTPError{}
	_ Error = &MultiError{}
)

// MultiError represents a slice of errors
type MultiError []Error

// Error returns the error message
func (e MultiError) Error() string {
	msg := []string{}
	for _, err := range e {
		msg = append(msg, err.Error())
	}

	return strings.Join(msg, ";")
}

// Prepare prepares the error for rendering
func (e MultiError) Prepare() Error {
	merr := MultiError{}
	for _, err := range e {
		merr = append(merr, err.Prepare())
	}
	return merr
}

// Fields returns all fields that should be logged
func (e MultiError) Fields() log.Fields {
	fields := log.Fields{}

	for index, err := range e {
		key := fmt.Sprintf("errors[%d]", index)
		fields[key] = reason(err)
	}

	return fields
}

// HTTPError is a more feature rich implementation of error interface inspired
// by PostgreSQL error style guide
type HTTPError struct {
	Code    int               `json:"code,omitempty" xml:"code,omitempty"`
	Message string            `json:"message" xml:"message"`
	Reason  error             `json:"reason,omitempty" xml:"reason,omitempty"`
	Details []string          `json:"details,omitempty" xml:"details,omitempty"`
	Stack   errors.StackTrace `json:"-" xml:"-"`
}

// NewError returns an error with error code and error messages provided in
// function params
func NewError(code int, msg string, details ...string) *HTTPError {
	if code <= 0 {
		code = CodeInternal
	}

	return &HTTPError{
		Code:    code,
		Message: msg,
		Details: details,
		Stack:   NewStack().StackTrace(),
	}
}

// Error returns the error message
func (e *HTTPError) Error() string {
	return e.Message
}

// StackTrace returns the stack trace
func (e *HTTPError) StackTrace() errors.StackTrace {
	return e.Stack
}

// Fields returns the fields that should be logged
func (e *HTTPError) Fields() log.Fields {
	fields := log.Fields{}

	if e.Code > 0 {
		fields["code"] = e.Code
	}

	if e.Reason != nil {
		fields["reason"] = reason(e.Reason)
	}

	for index, msg := range e.Details {
		key := fmt.Sprintf("details[%d]", index)
		fields[key] = msg
	}

	return fields
}

// Wrap wraps the actual error
func (e *HTTPError) Wrap(err error) *HTTPError {
	e.Reason = err
	e.Stack = NewStack().StackTrace()
	return e
}

// Cause returns the real reason for the error
func (e *HTTPError) Cause() error {
	if e.Reason == nil {
		return e
	}

	if reason, ok := e.Reason.(*HTTPError); ok {
		return reason.Cause()
	}

	return e.Reason
}

// Prepare prepares the error for rendering
func (e *HTTPError) Prepare() Error {
	err := &HTTPError{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}

	if e.Reason == nil {
		return err
	}

	if perr, ok := e.Reason.(Error); ok {
		err.Reason = perr.Prepare()
	} else {
		err.Reason = &HTTPError{Message: e.Reason.Error()}
	}

	return err
}

// FieldsFormatter are the error log fields
type FieldsFormatter log.Fields

// String returns the fields as string
func (f FieldsFormatter) String() string {
	buffer := &bytes.Buffer{}
	fields := log.Fields(f)

	for index, name := range fields.Names() {
		if index > 0 {
			fmt.Fprint(buffer, " ")
		}
		fmt.Fprintf(buffer, "%v:%v", name, fields.Get(name))
	}

	return fmt.Sprintf("[%s]", buffer.String())
}

func reason(err error) interface{} {
	switch errx := err.(type) {
	case *HTTPError:
		return FieldsFormatter(errx.Fields()).Add("message", errx.Message)
	case MultiError:
		return FieldsFormatter(errx.Fields())
	default:
		return errx.Error()
	}
}

// Add adds key to the formatter
func (f FieldsFormatter) Add(key string, value interface{}) FieldsFormatter {
	fields := log.Fields(f)
	fields[key] = value
	return f
}
