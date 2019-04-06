package rest

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validationFuncMap = make(map[string]validator.Func)
)

// RegisterValidation adds a validation with the given tag
func RegisterValidation(tag string, fn validator.Func) {
	validationFuncMap[tag] = fn
}

// Validate validates a data
func Validate(r *http.Request, data interface{}) error {
	v := validator.New()

	for key, fn := range validationFuncMap {
		if err := v.RegisterValidation(key, fn); err != nil {
			return errors.WrapSkipFrames(err, "validate", 2).AddTag("status", http.StatusInternalServerError)
		}
	}

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		switch render.GetRequestContentType(r) {
		case render.ContentTypeJSON:
			return tagName(field, "json")
		case render.ContentTypeXML:
			return tagName(field, "xml")
		case render.ContentTypeForm:
			return tagName(field, "form")
		default:
			return field.Name
		}
	})

	if err := v.StructCtx(r.Context(), data); err != nil {
		return errors.WrapSkipFrames(err, "validate", 2).AddTag("status", http.StatusUnprocessableEntity)
	}

	return nil
}

func tagName(field reflect.StructField, attr string) string {
	tag := field.Tag.Get(attr)

	if idx := strings.Index(tag, ","); idx != -1 {
		tag = tag[:idx]
	}

	if tag == "-" {
		tag = field.Name
	}

	return tag
}
