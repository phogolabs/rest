package httpr

import (
	"net/http"

	"github.com/go-chi/render"
)

// Render responds with success
func Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	if data == nil {
		return
	}

	if err, ok := data.(error); ok {
		respondErr(w, r, err)
	} else {
		respond(w, r, data)
	}
}

func respond(w http.ResponseWriter, r *http.Request, data interface{}) {
	response, ok := data.(*Response)

	if !ok {
		response = &Response{
			StatusCode: http.StatusOK,
			Data:       data,
		}
	}

	_ = render.Render(w, r, response)
}

func respondErr(w http.ResponseWriter, r *http.Request, err error) {
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

	response.Render(w, r)
}
