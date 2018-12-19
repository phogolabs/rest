package rest

import (
	"github.com/go-chi/render"
	"github.com/phogolabs/rest/middleware"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validationFuncMap = make(map[string]validator.Func)

	// GetLogger returns the associated request logger
	GetLogger = middleware.GetLogger
)

func init() {
	render.Decode = DefaultDecoder
}

// RegisterValidation adds a validation with the given tag
func RegisterValidation(tag string, fn validator.Func) {
	validationFuncMap[tag] = fn
}
