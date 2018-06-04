package httpr_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/http/httpr"
)

var _ = Describe("Response", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("Content-Type", "application/json")
	})

	It("respond successfully", func() {
		response := &httpr.Response{StatusCode: http.StatusCreated}
		httpr.Render(w, r, response)
		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		Expect(ok).To(BeTrue())
		Expect(status).To(Equal(http.StatusCreated))
	})

	It("sets the status code successfully", func() {
		response := &httpr.Response{StatusCode: http.StatusCreated}
		Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		Expect(ok).To(BeTrue())
		Expect(status).To(Equal(http.StatusCreated))
	})

	Context("when the output is generic type", func() {
		It("sets the status code successfully", func() {
			httpr.Render(w, r, "hello")
			status, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(status).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("hello"))
		})
	})

	Context("when the status code is not provided", func() {
		It("sets the status code 200 successfully", func() {
			response := &httpr.Response{}
			Expect(response.StatusCode).To(BeZero())
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			status, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(status).To(Equal(http.StatusOK))
		})
	})

	Context("when the kind is set", func() {
		It("does not set any kind", func() {
			response := &httpr.Response{Data: time.Now()}
			response.Meta.Kind = "test"
			Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
			Expect(response.Meta.Kind).To(Equal("test"))
		})
	})

	Context("when the kind is not set", func() {
		Context("when the data is nil", func() {
			It("does not set any kind", func() {
				response := &httpr.Response{}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(BeEmpty())
			})
		})

		Context("when the data is struct", func() {
			It("sets the kind successfully", func() {
				response := &httpr.Response{Data: time.Now()}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(Equal("time"))
			})
		})

		Context("when the data is not struct", func() {
			It("does not set the kind successfully", func() {
				response := &httpr.Response{Data: 5}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(Equal(""))
			})
		})

		Context("when the data is slice of struct", func() {
			It("sets the kind successfully", func() {
				arr := []*time.Time{}
				response := &httpr.Response{Data: &arr}
				Expect(response.Render(httptest.NewRecorder(), r)).To(Succeed())
				Expect(response.Meta.Kind).To(Equal("time"))
			})
		})
	})
})
