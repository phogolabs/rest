package rho

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/gosuri/uitable"
	"github.com/lib/pq"
)

const (
	// ErrCodeParamRequired is an error code returned when the parameter is missing
	ErrCodeParamRequired = 20101
	// ErrCodeParamInvalid is an error code returned when the parameter's value is an invalid
	ErrCodeParamInvalid = 20102
	// ErrCodeQueryParamRequired is an error code returned when the query parameter is missing
	ErrCodeQueryParamRequired = 20103
	// ErrCodeQueryParamInvalid is an error code returned when the query parameter's value is an invalid
	ErrCodeQueryParamInvalid = 20104
	// ErrCodeUnknown when the error is unknown
	ErrCodeUnknown = 40000
	// ErrCodeConflict when the API request cannot be completed because the requested operation would conflict with an existing item.
	ErrCodeConflict = 40101
	// ErrCodeDuplicate when the requested operation failed because it tried to create a resource that already exists.
	ErrCodeDuplicate = 40102
	// ErrCodeDeleted when the request failed because the resource associated with the request has been deleted 410
	ErrCodeDeleted = 40103
	// ErrCodeConditionNotMet when the condition set in the request's was not met 416
	ErrCodeConditionNotMet = 40104
	// ErrCodeRequestedRangeNotSatisfiable when the request specified a range that cannot be satisfied 428.
	ErrCodeRequestedRangeNotSatisfiable = 40105
	// ErrCodeInternalError when the request failed due to an internal error 500.
	ErrCodeInternalError = 40106
	// ErrBackendError when a backend error occurred 503. Usually database.
	ErrBackendError = 40106
	// ErrBackendNotConnected when the request failed due to a connection error.
	ErrBackendNotConnected = 40107
	// ErrNotReady code when the API server is not ready to accept requests.
	ErrNotReady = 40108
	// ErrCodeValidationError when the one or more field are validated
	ErrCodeValidationError = 40109
)

var _ error = &Error{}

// Error is a more feature rich implementation of error interface inspired
// by PostgreSQL error style guide
type Error struct {
	Code    int      `json:"code" xml:"code"`
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

	switch errx := err.(type) {
	case pq.Error:
		response = PostgreSQLErrorResponse(errx)
	case *ErrorResponse:
		response = errx
	}

	if err := render.Render(w, r, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PostgreSQLErrorResponse creates a ErrorResponse for given PostgreSQL error
func PostgreSQLErrorResponse(err pq.Error) *ErrorResponse {
	return nil
}
