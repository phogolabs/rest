package httperr_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/phogolabs/rho/httperr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the err is not recognized", func() {
		It("handles the error correctly", func() {
			errx := httperr.JSONError(fmt.Errorf("Oh no!"))
			Expect(errx.StatusCode).To(Equal(http.StatusInternalServerError))
			Expect(errx.Err.Code).To(Equal(httperr.CodeInternal))
			Expect(errx.Err.Message).To(Equal("JSON Error"))
		})
	})

	Context("when the err is json.InvalidUnmarshalError", func() {
		It("handles the error correctly", func() {
			err := &json.InvalidUnmarshalError{Type: reflect.TypeOf(r)}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal json body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is json.UnmarshalFieldError", func() {
		It("handles the error correctly", func() {
			typ := reflect.TypeOf(*r)
			err := &json.UnmarshalFieldError{Type: typ, Field: typ.Field(0), Key: "StatusCode"}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal json body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is UnmarshalTypeError", func() {
		It("handles the error correctly", func() {
			err := &json.UnmarshalTypeError{Type: reflect.TypeOf(*r)}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal json body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is UnsupportedValueError", func() {
		It("handles the error correctly", func() {
			err := &json.UnsupportedValueError{Value: reflect.ValueOf(*r)}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is UnsupportedTypeError", func() {
		It("handles the error correctly", func() {
			err := &json.UnsupportedTypeError{Type: reflect.TypeOf(*r)}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is InvalidUTF8Error", func() {
		It("handles the error correctly", func() {
			err := &json.InvalidUTF8Error{S: "Oh no!"}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is MarshalerError", func() {
		It("handles the error correctly", func() {
			err := &json.MarshalerError{Err: fmt.Errorf("Oh no!"), Type: reflect.TypeOf(*r)}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})
