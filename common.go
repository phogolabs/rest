package rest

import (
	"github.com/go-chi/render"

	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	render.Decode = DefaultDecoder
	render.Respond = DefaultResponder
}

// RegisterValidation adds a validation with the given tag
func RegisterValidation(tag string, fn validator.Func) {
	validationFuncMap[tag] = fn
}
