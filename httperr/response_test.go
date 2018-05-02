package httperr_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httperr"
)

var _ = Describe("ErrResponse", func() {
	var r *http.Request

	BeforeEach(func() {
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	It("returns the underlying error message", func() {
		err := httperr.New(12345, "Oh no!")

		response := &httperr.Response{Err: err}
		Expect(response.Error()).To(Equal(err.Message))
	})

	It("set the status code", func() {
		leaf := fmt.Errorf("Inner Error")

		err := httperr.New(201, "Oh no!", "Unexpected error")
		err.Wrap(httperr.New(202, "Madness").Wrap(leaf))

		response := &httperr.Response{StatusCode: http.StatusForbidden, Err: err}
		Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		Expect(ok).To(BeTrue())
		Expect(status).To(Equal(http.StatusForbidden))
	})

	Context("when the status code is not provided", func() {
		It("set the status code internal server error", func() {
			err := httperr.New(12345, "Oh no!")
			response := &httperr.Response{Err: err}
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			status, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(status).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when the error code is not set", func() {
		It("set the unknown erro code", func() {
			err := httperr.New(0, "Oh no!")
			response := &httperr.Response{Err: err}
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			Expect(response.Err.Code).To(Equal(httperr.ErrUnknown))
		})
	})
})

var _ = Describe("RespondErr", func() {
	var r *http.Request

	BeforeEach(func() {
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	It("respond with the provided information", func() {
		w := httptest.NewRecorder()
		err := httperr.New(2000, "Oh no!")
		httperr.Respond(w, r, http.StatusForbidden, err)

		Expect(w.Code).To(Equal(http.StatusForbidden))
		Expect(w.Body.String()).To(ContainSubstring(`{"error":{"code":2000,"message":"Oh no!"}}`))
	})
})
