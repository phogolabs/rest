package rho

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/goware/errorx"
	"github.com/lib/pq"
)

func init() {
	errorx.SetVerbosity(1)
}

const (
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

var (
	_ render.Renderer = &ErrorResponse{}
	_ error           = &ErrorResponse{}
)

// Error is a more feature rich implementation of error interface inspired
// by PostgreSQL error style guide
type Error = errorx.Errorx

// ErrorResponse represents a HTTP error response
type ErrorResponse struct {
	StatusCode int
	Err        *Error
}

// Error returns the error message from the underlying error
func (e *ErrorResponse) Error() string {
	return e.Err.Error()
}

// Render renders a single error and respond to the client request.
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if e.StatusCode <= 0 {
		return fmt.Errorf("Invalid status code: %d", e.StatusCode)
	}

	render.Status(r, e.StatusCode)
	return nil
}

// MarshalJSON encodes the error into JSON
func (e *ErrorResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error)
}

// MarshalXML encodes the error into XML
func (e *ErrorResponse) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.EncodeElement(e.Error, start)
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
