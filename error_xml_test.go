package rho_test

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho"
)

var _ = Describe("ErrorXml", func() {
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
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.SyntaxError", func() {
		It("handles the error correctly", func() {
			err := &xml.SyntaxError{Msg: "oh no!"}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.TagPathError", func() {
		It("handles the error correctly", func() {
			err := &xml.TagPathError{Struct: reflect.TypeOf(*r)}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.UnsupportedTypeError", func() {
		It("handles the error correctly", func() {
			err := &xml.UnsupportedTypeError{Type: reflect.TypeOf(*r)}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal xml"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})
