package httpr

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
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

// ErrorRenderer renders error as HTTP response
type ErrorRenderer interface {
	render.Renderer
	error
}

// LoggableError logs an error by providing extra information
type LoggableError interface {
	Fields() log.Fields
	error
}

var (
	_ ErrorRenderer = &Error{}
	_ LoggableError = &Error{}
	_ LoggableError = &ErrorList{}
)

// ErrorList represents a slice of errors
type ErrorList []error

// Error returns the error message
func (e ErrorList) Error() string {
	msg := []string{}
	for _, err := range e {
		msg = append(msg, err.Error())
	}

	return strings.Join(msg, ";")
}

// Fields returns all fields that should be logged
func (e ErrorList) Fields() log.Fields {
	fields := log.Fields{}

	for index, err := range e {
		key := fmt.Sprintf("errors[%d]", index)
		fields[key] = reason(err)
	}

	return fields
}

// Error is a more feature rich implementation of error interface inspired
// by PostgreSQL error style guide
type Error struct {
	Code    int               `json:"code,omitempty" xml:"code,omitempty"`
	Message string            `json:"message" xml:"message"`
	Reason  error             `json:"reason,omitempty" xml:"reason,omitempty"`
	Details []string          `json:"details,omitempty" xml:"details,omitempty"`
	Status  int               `json:"-" xml:"-"`
	Stack   errors.StackTrace `json:"-" xml:"-"`
}

// NewError returns an error with error code and error messages provided in
// function params
func NewError(code int, msg string, details ...string) *Error {
	if code <= 0 {
		code = CodeInternal
	}

	return &Error{
		Status:  http.StatusInternalServerError,
		Code:    code,
		Message: msg,
		Details: details,
		Stack:   NewStack().StackTrace(),
	}
}

// Error returns the error message
func (e *Error) Error() string {
	return e.Message
}

// StackTrace returns the stack trace
func (e *Error) StackTrace() errors.StackTrace {
	return e.Stack
}

// Fields returns the fields that should be logged
func (e *Error) Fields() log.Fields {
	fields := log.Fields{}

	if e.Status > 0 {
		fields["status"] = e.Status
	}

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
func (e *Error) Wrap(err error) *Error {
	e.Reason = err
	e.Stack = NewStack().StackTrace()
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

// WithStatus creates a new error with status
func (e *Error) WithStatus(status int) *Error {
	e.Status = status
	return e
}

// Render renders a single error and respond to the client request.
func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
	if logEntry := middleware.GetLogEntry(r); logEntry != nil {
		logEntry.Panic(e, nil)
	}

	response := render.M{
		"error": prepare(e),
	}

	render.Status(r, e.Status)
	render.Respond(w, r, response)
	return nil
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

// Add adds key to the formatter
func (f FieldsFormatter) Add(key string, value interface{}) FieldsFormatter {
	fields := log.Fields(f)
	fields[key] = value
	return f
}

// ConvError creates a Response for given strconv error
func ConvError(err error) *Error {
	errx := NewError(CodeInvalid, "Unable to parse number")
	return errx.WithStatus(http.StatusUnprocessableEntity).Wrap(err)
}

// TimeError creates a Response for given time error
func TimeError(err error) *Error {
	errx := NewError(CodeInvalid, "Unable to parse date time")
	return errx.WithStatus(http.StatusUnprocessableEntity).Wrap(err)
}

// JSONError creates a Response for given json error
func JSONError(err error) *Error {
	const (
		jsonMsg          = "JSON Error"
		jsonUnmarshalMsg = "Unable to unmarshal json body"
		jsonMarshalMsg   = "Unable to marshal json"
	)

	switch err.(type) {
	case *json.InvalidUnmarshalError:
		errx := NewError(CodeInternal, jsonUnmarshalMsg)
		return errx.WithStatus(http.StatusInternalServerError).Wrap(err)
	case *json.UnmarshalFieldError, *json.UnmarshalTypeError:
		errx := NewError(CodeInvalid, jsonUnmarshalMsg)
		return errx.WithStatus(http.StatusBadRequest).Wrap(err)
	case *json.UnsupportedTypeError, *json.UnsupportedValueError, *json.InvalidUTF8Error, *json.MarshalerError:
		errx := NewError(CodeInternal, jsonMarshalMsg)
		return errx.WithStatus(http.StatusInternalServerError).Wrap(err)
	default:
		errx := NewError(CodeInternal, jsonMsg)
		return errx.WithStatus(http.StatusInternalServerError).Wrap(err)
	}
}

// XMLError creates a Response for given xml error
func XMLError(err error) *Error {
	const (
		xmlMsg          = "XML Error"
		xmlUnmarshalMsg = "Unable to unmarshal xml body"
		xmlMarshalMsg   = "Unable to marshal xml"
	)

	switch err.(type) {
	case xml.UnmarshalError, *xml.SyntaxError, *xml.TagPathError:
		errx := NewError(CodeInvalid, xmlUnmarshalMsg)
		return errx.WithStatus(http.StatusBadRequest).Wrap(err)
	case *xml.UnsupportedTypeError:
		errx := NewError(CodeInternal, xmlMarshalMsg)
		return errx.WithStatus(http.StatusInternalServerError).Wrap(err)
	default:
		errx := NewError(CodeInternal, xmlMsg)
		return errx.WithStatus(http.StatusInternalServerError).Wrap(err)
	}
}

// UnknownError handles httperr
func UnknownError(err error) *Error {
	switch errx := err.(type) {
	case *Error:
		return errx
	default:
		unerr := NewError(CodeInternal, "Internal Error")
		return unerr.WithStatus(http.StatusInternalServerError).Wrap(err)
	}
}

const (
	// Class 08 - Connection Exception
	pgConnClassErr = "08"
	// Class 22 - Data Exception
	pgDataClassErr = "22"
	// Class 23 - Integrity Constraint Violation
	pgContraintClassErr = "23"
	// Class 57 - Operator Intervention
	pgOpIntClassErr = "57"
)

// PGError creates a Response for given PostgreSQL error
func PGError(err error) *Error {
	var (
		pgErr    = err.(pq.Error)
		response *Error
	)

	switch pgErr.Code[:2] {
	case pgConnClassErr:
		response = NewError(CodeBackendNotConnected, "Connection Error")
		response = response.WithStatus(http.StatusInternalServerError)
		response = response.Wrap(err)
	case pgDataClassErr:
		response = PGDataError(pgErr)
	case pgContraintClassErr:
		response = PGIntegrityError(pgErr)
	case pgOpIntClassErr:
		response = NewError(CodeBackendNotReady, "Operator Intervention")
		response = response.WithStatus(http.StatusInternalServerError)
		response = response.Wrap(err)
	default:
		response = NewError(CodeBackend, "Database Error")
		response = response.WithStatus(http.StatusInternalServerError)
		response = response.Wrap(err)
	}

	return response
}

// PGIntegrityError handles PG integrity errors
func PGIntegrityError(err pq.Error) *Error {
	errx := NewError(CodeConflict, "Integrity Constraint Violation")

	switch err.Code {
	// "23505": "unique_violation",
	case "23505":
		errx.Code = CodeDuplicate
	// "23514": "check_violation"
	// "23P01": "exclusion_violation"
	case "23514", "23P01":
		errx.Code = CodeConditionNotMet
	}

	errx = errx.WithStatus(http.StatusConflict)
	errx = errx.Wrap(err)

	return errx
}

// PGDataError handles PG integrity errors
func PGDataError(err pq.Error) *Error {
	errx := NewError(CodeConflict, "Data Error")

	switch err.Code {
	// "22003": "numeric_value_out_of_range",
	// "22008": "datetime_field_overflow",
	// "22015": "interval_field_overflow",
	// "22022": "indicator_overflow",
	// "22P01": "floating_point_exception",
	case "22003", "22008", "22015", "22022", "22P01":
		errx.Code = CodeOutOfrange
	// "22004": "null_value_not_allowed",
	// "22002": "null_value_no_indicator_parameter",
	case "22002", "22004":
		errx.Code = CodeConditionNotMet
	}

	errx = errx.WithStatus(http.StatusUnprocessableEntity)
	errx = errx.Wrap(err)

	return errx
}

// ValidationError is an error which occurrs duering validation
func ValidationError(err error) *Error {
	rerrx := NewError(CodeConditionNotMet, "validation failed")
	rerrx = rerrx.WithStatus(http.StatusUnprocessableEntity)

	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		return rerrx.Wrap(err)
	}

	errs := ErrorList{}

	for _, ferr := range errors {
		if err, ok := ferr.(error); ok {
			errs = append(errs, err)
		}
	}

	return rerrx.Wrap(errs)
}

func reason(err error) interface{} {
	switch errx := err.(type) {
	case *Error:
		return FieldsFormatter(errx.Fields()).Add("message", errx.Message)
	case ErrorList:
		return FieldsFormatter(errx.Fields())
	default:
		return errx.Error()
	}
}

func prepare(err error) error {
	if err == nil {
		return err
	}

	switch errx := err.(type) {
	case *Error:
		result := *errx
		result.Reason = prepare(result.Reason)
		return &result
	case ErrorList:
		result := ErrorList{}
		for _, item := range errx {
			result = append(result, prepare(item))
		}
		return result
	default:
		return &Error{Message: err.Error()}
	}
}
