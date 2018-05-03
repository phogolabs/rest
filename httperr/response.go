package httperr

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
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
	return e.Err.Error()
}

// Render renders a single error and respond to the client request.
func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if e.StatusCode <= 0 {
		e.StatusCode = http.StatusInternalServerError
	}

	e.Err = e.Err.prepare()

	if logEntry := middleware.GetLogEntry(r); logEntry != nil {
		logEntry.Panic(e.Err, nil)
	}

	render.Status(r, e.StatusCode)
	return nil
}

// Respond handles the error based on pre-defined list of responses
func Respond(w http.ResponseWriter, r *http.Request, err error) {
	var response *Response

	switch pkgName(err) {
	case "github.com/lib/pq":
		response = PGError(err)
	case "github.com/phogolabs/rho/httperr":
		response = HTTPError(err)
	case "encoding/json":
		response = JSONError(err)
	case "encoding/xml":
		response = XMLError(err)
	case "strconv":
		response = ConvError(err)
	case "time":
		response = TimeError(err)
	case "gopkg.in/go-playground/validator.v9":
		response = ValidationError(err)
	default:
		response = New(CodeInternal, "Internal Error").With(http.StatusInternalServerError)
	}

	prepare(response, err)

	// Response never fails
	_ = render.Render(w, r, response)
}

func prepare(response *Response, err error) {
	if err != response && err != response.Err && response.Err.Reason == nil {
		response.Err.Reason = err
		response.Err.Stack = NewStack().StackTrace()
	}
}

// HTTPError handles httperr
func HTTPError(err error) *Response {
	switch errx := err.(type) {
	case *Response:
		return errx
	case *Error:
		return errx.With(http.StatusInternalServerError)
	default:
		return New(CodeInternal, "Internal Error").With(http.StatusInternalServerError)
	}
}
