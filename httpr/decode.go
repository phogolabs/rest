package httpr

import (
	"net/http"

	"github.com/go-chi/render"
)

// Bind binds a request into a struct
func Bind(r *http.Request, binder render.Binder) error {
	if err := render.Bind(r, binder); err != nil {
		return err
	}

	if err := DefaultValidator.Struct(binder); err != nil {
		return err
	}

	return nil
}
