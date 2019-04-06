package rest

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/errors"
	"github.com/go-playground/form"
)

// Decode is a package-level variable set to our default Decoder. We do this
// because it allows you to set render.Decode to another function with the
// same function signature, while also utilizing the render.Decoder() function
// itself. Effectively, allowing you to easily add your own logic to the package
// defaults. For example, maybe you want to impose a limit on the number of
// bytes allowed to be read from the request body.
func Decode(r *http.Request, v interface{}) error {
	var err error

	switch render.GetRequestContentType(r) {
	case render.ContentTypeJSON:
		err = render.DecodeJSON(r.Body, v)
	case render.ContentTypeXML:
		err = render.DecodeXML(r.Body, v)
	case render.ContentTypeForm:
		err = DecodeForm(r, v)
	default:
		err = errors.New("render: unable to automatically decode the request content type")
	}

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		err = errors.WrapSkipFrames(err, "decode", 2).AddTag("status", http.StatusBadRequest)
	}

	if err == nil {
		err = Validate(r, v)
	}

	return err
}

// DecodeForm decodes an entity from form fields
func DecodeForm(r *http.Request, v interface{}) (err error) {
	decoder := form.NewDecoder()

	if err = r.ParseForm(); err == nil {
		err = decoder.Decode(v, r.Form)
	}

	return
}

// DecodePath decodes an entity from path
func DecodePath(r *http.Request, v interface{}) error {
	decoder := form.NewDecoder()
	decoder.SetTagName("path")

	var (
		values = url.Values{}
		ctx    = chi.RouteContext(r.Context())
	)

	for index, key := range ctx.URLParams.Keys {
		values.Add(key, ctx.URLParams.Values[index])
	}

	return decoder.Decode(v, values)
}

// DecodeQuery decodes an entity from query
func DecodeQuery(r *http.Request, v interface{}) error {
	decoder := form.NewDecoder()
	decoder.SetTagName("query")

	values := url.Values{}

	if r.URL != nil {
		values = r.URL.Query()
	}

	return decoder.Decode(v, values)
}

// DecodeHeader decodes an entity from query
func DecodeHeader(r *http.Request, v interface{}) error {
	decoder := form.NewDecoder()
	decoder.SetTagName("header")

	values := url.Values(r.Header)
	return decoder.Decode(v, values)
}
