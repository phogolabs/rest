package httpr

import (
	"net/http"

	"github.com/go-chi/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// DefaultValidator is the default payload validator
var DefaultValidator = validator.New()

var _ render.Renderer = &Response{}

// ResponseMeta keeps meta information for the successful response
type ResponseMeta struct {
	// Kind of the data kept
	Kind string `json:"kind,omitempty" xml:"kind,attr"`
}

// Response represents a payload of successful response
type Response struct {
	// StatusCode for the response. Default 200 OK
	StatusCode int `json:"-" xml:"-"`
	// Meta is the metadata for this response
	Meta ResponseMeta `json:"meta,omitempty" xml:"meta,omitempty"`
	// Date of this response
	Data interface{} `json:"data,omitempty" xml:"data,omitempty"`
}

// Render renders a single response and respond to the client request.
func (p *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if p.StatusCode <= 0 {
		p.StatusCode = http.StatusOK
	}

	if p.Meta.Kind == "" {
		p.Meta.Kind = typeName(p.Data)
	}

	render.Status(r, p.StatusCode)
	return nil
}

// Render responds with success
func Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	if data == nil {
		return
	}

	response, ok := data.(*Response)

	if !ok {
		response = &Response{
			StatusCode: http.StatusOK,
			Data:       data,
		}
	}

	_ = render.Render(w, r, response)
}

// RenderError responds with error to the client
func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	var response render.Renderer

	switch pkgName(err) {
	case "github.com/lib/pq":
		response = PGError(err)
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
		response = UnknownError(err)
	}

	if response != nil {
		response.Render(w, r)
	}
}
