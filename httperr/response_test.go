package httperr_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httperr"
)

var _ = Describe("Response", func() {
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

	Context("when the err is package err", func() {
		It("handles the error correctly", func() {
			err := httperr.New(201, "Oh no!", "Unexpected Error")
			errx := httperr.HTTPError(err)
			Expect(errx.StatusCode).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when the err is not recognized", func() {
		It("handles the error correctly", func() {
			err := httperr.HTTPError(nil)
			Expect(err).NotTo(BeNil())
			Expect(err.StatusCode).To(Equal(http.StatusInternalServerError))
			Expect(err.Err.Code).To(Equal(httperr.CodeInternal))
			Expect(err.Err.Message).To(Equal("Internal Error"))
		})
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
			Expect(response.Err.Code).To(Equal(httperr.CodeInternal))
		})
	})
})

var _ = Describe("RespondErr", func() {
	var (
		r      *http.Request
		w      *httptest.ResponseRecorder
		buffer *bytes.Buffer
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		buffer = &bytes.Buffer{}

		formatter := &middleware.DefaultLogFormatter{
			Logger: log.New(buffer, "", log.LstdFlags),
		}

		r = httptest.NewRequest("GET", "http://example.com", nil)
		r = middleware.WithLogEntry(r, formatter.NewLogEntry(r))
	})

	It("respond with the provided information", func() {
		err := httperr.New(2000, "Oh no!", "Inner Error")
		httperr.Respond(w, r, err.With(http.StatusForbidden))

		Expect(w.Code).To(Equal(http.StatusForbidden))
		Expect(w.Body.String()).To(ContainSubstring(`{"error":{"code":2000,"message":"Oh no!","details":["Inner Error"]}}`))
		Expect(buffer.String()).To(ContainSubstring("example"))
	})

	Context("when the error is nil", func() {
		It("does not respond with the provided information", func() {
			Expect(func() { httperr.Respond(w, r, nil) }).NotTo(Panic())
		})
	})

	Context("when the error is slice", func() {
		It("does not respond with the provided information", func() {
			Expect(func() { httperr.Respond(w, r, FakeSliceErr{}) }).NotTo(Panic())
		})
	})
})
