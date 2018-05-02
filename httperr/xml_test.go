package httperr_test

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httperr"
)

var _ = Describe("XML Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the err is xml.UnmarshalError", func() {
		It("handles the error correctly", func() {
			err := xml.UnmarshalError("oh no")
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.SyntaxError", func() {
		It("handles the error correctly", func() {
			err := &xml.SyntaxError{Msg: "oh no!"}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.TagPathError", func() {
		It("handles the error correctly", func() {
			err := &xml.TagPathError{Struct: reflect.TypeOf(*r)}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.UnsupportedTypeError", func() {
		It("handles the error correctly", func() {
			err := &xml.UnsupportedTypeError{Type: reflect.TypeOf(*r)}
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal xml"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})
