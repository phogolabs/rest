package rho

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/gosuri/uitable"
)

const (
	// ErrParamRequired is an error code returned when the parameter is missing
	ErrParamRequired = 20101
	// ErrParamInvalid is an error code returned when the parameter's value is an invalid
	ErrParamInvalid = 20102
)

const (
	// ErrUnknown when the error is unknown
	ErrUnknown = 40000
	// ErrConflict when the API request cannot be completed because the requested operation would conflict with an existing item.
	ErrConflict = 40101
	// ErrDuplicate when the requested operation failed because it tried to create a resource that already exists.
	ErrDuplicate = 40102
	// ErrDeleted when the request failed because the resource associated with the request has been deleted 410
	ErrDeleted = 40103
	// ErrConditionNotMet when the condition set in the request's was not met 416
	ErrConditionNotMet = 40104
	// ErrOutOfrange when the request specified a range that cannot be satisfied 428.
	ErrOutOfrange = 40105
)

const (
	// ErrInvalid when the provided payload is invalid
	ErrInvalid = 40222
	// ErrBackend when a backend error occurred 503. Usually database.
	ErrBackend = 40106
	// ErrBackendNotConnected when the request failed due to a connection error.
	ErrBackendNotConnected = 40107
	// ErrNotReady code when the API server is not ready to accept requests.
	ErrNotReady = 40108
	// ErrInternal when the request failed due to an internal error 500.
	ErrInternal = 40106
	// ErrValidationError when the one or more field are validated
	ErrValidationError = 40109
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

// NewError returns an error with error code and error messages provided in
// function params
func NewError(code int, msg ...string) *Error {
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

// RespondErr responses with error
func RespondErr(w http.ResponseWriter, r *http.Request, statusCode int, err *Error) {
	response := &ErrorResponse{
		StatusCode: statusCode,
		Err:        err,
	}

	if err := render.Render(w, r, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleErr handles the error based on pre-defined list of responses
func HandleErr(w http.ResponseWriter, r *http.Request, err error) {
	var response *ErrorResponse

	switch pkgName(err) {
	case "github.com/lib/pq":
		response = PGError(err)
	case "github.com/phogolabs/rho":
		response = err.(*ErrorResponse)
	case "encoding/json":
		response = JSONError(err)
	case "encoding/xml":
		response = XMLError(err)
	case "strconv":
		response = ConvError(err)
	case "time":
		response = TimeError(err)
	default:
		response = &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(ErrUnknown, "Unknown Error").Wrap(err),
		}
	}

	if err := render.Render(w, r, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
