package rho_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho"
)

var _ = Describe("Response", func() {
	var r *http.Request

	BeforeEach(func() {
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	It("sets the status code successfully", func() {
		response := &rho.Response{StatusCode: http.StatusCreated}
		Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		Expect(ok).To(BeTrue())
		Expect(status).To(Equal(http.StatusCreated))
	})

	Context("when the status code is not provided", func() {
		It("sets the status code 200 successfully", func() {
			response := &rho.Response{}
			Expect(response.StatusCode).To(BeZero())
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			status, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(status).To(Equal(http.StatusOK))
		})
	})

	Context("when the kind is set", func() {
		It("does not set any kind", func() {
			response := &rho.Response{Data: time.Now()}
			response.Meta.Kind = "test"
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			Expect(response.Meta.Kind).To(Equal("test"))
		})
	})

	Context("when the kind is not set", func() {
		Context("when the data is nil", func() {
			It("does not set any kind", func() {
				response := &rho.Response{}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(BeEmpty())
			})
		})

		Context("when the data is struct", func() {
			It("sets the kind successfully", func() {
				response := &rho.Response{Data: time.Now()}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(Equal("time"))
			})
		})

		Context("when the data is not struct", func() {
			It("does not set the kind successfully", func() {
				response := &rho.Response{Data: 5}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(Equal(""))
			})
		})

		Context("when the data is slice of struct", func() {
			It("sets the kind successfully", func() {
				arr := []*time.Time{}
				response := &rho.Response{Data: &arr}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(Equal("time"))
			})
		})
	})
})

var _ = Describe("ErrResponse", func() {
	var r *http.Request

	BeforeEach(func() {
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	It("returns the underlying error message", func() {
		err := rho.NewError(12345, "Oh no!")

		response := &rho.ErrorResponse{Err: err}
		Expect(response.Error()).To(Equal(err.Message))
	})

	It("set the status code", func() {
		leaf := fmt.Errorf("Inner Error")

		err := rho.NewError(201, "Oh no!", "Unexpected error")
		err.Wrap(rho.NewError(202, "Madness").Wrap(leaf))

		response := &rho.ErrorResponse{StatusCode: http.StatusForbidden, Err: err}
		Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		Expect(ok).To(BeTrue())
		Expect(status).To(Equal(http.StatusForbidden))
	})

	Context("when the status code is not provided", func() {
		It("set the status code internal server error", func() {
			err := rho.NewError(12345, "Oh no!")
			response := &rho.ErrorResponse{Err: err}
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			status, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(status).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when the error code is not set", func() {
		It("set the unknown erro code", func() {
			err := rho.NewError(0, "Oh no!")
			response := &rho.ErrorResponse{Err: err}
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			Expect(response.Err.Code).To(Equal(rho.ErrUnknown))
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
		err := rho.NewError(2000, "Oh no!")
		rho.RespondErr(w, r, http.StatusForbidden, err)

		Expect(w.Code).To(Equal(http.StatusForbidden))
		Expect(w.Body.String()).To(ContainSubstring(`{"error":{"code":2000,"message":"Oh no!"}}`))
	})
})
