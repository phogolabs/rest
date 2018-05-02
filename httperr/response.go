package httperr

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	_ render.Renderer = &Response{}
	_ error           = &Response{}
)

// Response represents a HTTP error response
type Response struct {
	StatusCode int    `json:"-" xml:"-"`
	Err        *Error `json:"error" xml:"error"`
}

// Error returns the error message from the underlying error
func (e *Response) Error() string {
	return e.Err.Message
}

// Render renders a single error and respond to the client request.
func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if e.StatusCode <= 0 {
		e.StatusCode = http.StatusInternalServerError
	}

	e.Err = e.Err.prepare()

	if e.Err.Code <= 0 {
		e.Err.Code = ErrUnknown
	}

	render.Status(r, e.StatusCode)
	return nil
}

// Respond responses with error
func Respond(w http.ResponseWriter, r *http.Request, statusCode int, err *Error) {
	response := &Response{
		StatusCode: statusCode,
		Err:        err,
	}

	if err := render.Render(w, r, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handle handles the error based on pre-defined list of responses
func Handle(w http.ResponseWriter, r *http.Request, err error) {
	var response *Response

	switch pkgName(err) {
	case "github.com/lib/pq":
		response = PGError(err)
	case "github.com/phogolabs/rho/httperr":
		response = err.(*Response)
	case "encoding/json":
		response = JSONError(err)
	case "encoding/xml":
		response = XMLError(err)
	case "strconv":
		response = ConvError(err)
	case "time":
		response = TimeError(err)
	default:
		response = &Response{
			StatusCode: http.StatusInternalServerError,
			Err:        New(ErrUnknown, "Unknown Error").Wrap(err),
		}
	}

	if err := render.Render(w, r, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
