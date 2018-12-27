package rest

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/form"
	validator "gopkg.in/go-playground/validator.v9"
)

type (
	// Binder interface for managing request payloads.
	Binder = render.Binder
)

// Bind decodes a request body and executes the Binder method of the
// payload structure.
func Bind(r *http.Request, v Binder) error {
	err := render.Bind(r, v)

	if err == nil {
		err = Validate(r, v)
	}

	return err
}

// Decode is a package-level variable set to our default Decoder. We do this
// because it allows you to set render.Decode to another function with the
// same function signature, while also utilizing the render.Decoder() function
// itself. Effectively, allowing you to easily add your own logic to the package
// defaults. For example, maybe you want to impose a limit on the number of
// bytes allowed to be read from the request body.
func Decode(r *http.Request, v interface{}) error {
	err := DefaultDecoder(r, v)

	if err == nil {
		err = Validate(r, v)
	}

	return err
}

// DefaultDecoder is the default decoder
func DefaultDecoder(r *http.Request, v interface{}) error {
	var err error

	switch render.GetRequestContentType(r) {
	case render.ContentTypeForm:
		decoder := form.NewDecoder()

		if err = r.ParseForm(); err != nil {
			return err
		}

		err = decoder.Decode(v, r.Form)
	default:
		err = render.DefaultDecoder(r, v)
	}

	return err
}

// Validate validates a data
func Validate(r *http.Request, data interface{}) error {
	v := validator.New()

	for key, fn := range validationFuncMap {
		if err := v.RegisterValidation(key, fn); err != nil {
			return err
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
		return err
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
