package rho_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Param", func() {
	var (
		r   *http.Request
		ctx *chi.Context
	)

	BeforeEach(func() {
		ctx = chi.NewRouteContext()
	})

	JustBeforeEach(func() {
		r = httptest.NewRequest("GET", "http://example.com", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	})

	Describe("URLParamUUID", func() {
		It("parses the values successfully", func() {
			id := uuid.NewV4()
			ctx.URLParams.Add("id", id.String())

			value, err := rho.URLParamUUID(r, "id")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(id))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := rho.URLParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*rho.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Parameter 'id' is required"))
				Expect(value).To(Equal(uuid.Nil))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				ctx.URLParams.Add("id", "wrong-uuid")

				value, err := rho.URLParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*rho.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Parameter 'id' is not valid UUID"))
				Expect(value).To(Equal(uuid.Nil))
			})
		})
	})
})
