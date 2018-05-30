package httpr_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/chi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/http/httpr"
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

			value, err := httpr.URLParamUUID(r, "id")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(id))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'id' is required"))
				Expect(value).To(Equal(uuid.Nil))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				ctx.URLParams.Add("id", "wrong-uuid")

				value, err := httpr.URLParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'id' is not valid UUID"))
				Expect(value).To(Equal(uuid.Nil))
			})

			It("returns a nil value", func() {
				ctx.URLParams.Add("id", "wrong-uuid")
				Expect(httpr.URLParamUUIDOrNil(r, "id")).To(Equal(uuid.Nil))
			})

			It("returns the provided value", func() {
				id := uuid.NewV4()
				ctx.URLParams.Add("id", "wrong-uuid")
				Expect(httpr.URLParamUUIDOrValue(r, "id", id)).To(Equal(id))
			})
		})
	})

	Describe("URLParamInt", func() {
		It("parses the values successfully", func() {
			num := int64(123)
			ctx.URLParams.Add("num", "123")

			value, err := httpr.URLParamInt(r, "num", 0, 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLParamInt(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'num' is required"))
				Expect(value).To(Equal(int64(0)))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				ctx.URLParams.Add("num", "number")

				value, err := httpr.URLParamInt(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'num' is not valid integer number"))
				Expect(value).To(Equal(int64(0)))
			})

			It("returns the provided value", func() {
				value := int64(200)
				ctx.URLParams.Add("num", "number")
				Expect(httpr.URLParamIntOrValue(r, "num", 0, 64, value)).To(Equal(value))
			})
		})
	})

	Describe("URLParamUint", func() {
		It("parses the values successfully", func() {
			num := uint64(123)
			ctx.URLParams.Add("num", "123")

			value, err := httpr.URLParamUint(r, "num", 0, 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is negative", func() {
			It("parses the values successfully", func() {
				ctx.URLParams.Add("num", "-123")

				value, err := httpr.URLParamUint(r, "num", 0, 64)
				Expect(err).NotTo(BeNil())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'num' is not valid unsigned integer number"))
				Expect(value).To(Equal(uint64(0)))
			})
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLParamUint(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'num' is required"))
				Expect(value).To(Equal(uint64(0)))
			})

			It("returns the provided value", func() {
				value := uint64(200)
				ctx.URLParams.Add("num", "number")
				Expect(httpr.URLParamUintOrValue(r, "num", 0, 64, value)).To(Equal(value))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				ctx.URLParams.Add("num", "number")

				value, err := httpr.URLParamUint(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'num' is not valid unsigned integer number"))
				Expect(value).To(Equal(uint64(0)))
			})
		})
	})

	Describe("URLParamFloat", func() {
		It("parses the values successfully", func() {
			num := float64(123.11)
			ctx.URLParams.Add("num", "123.11")

			value, err := httpr.URLParamFloat(r, "num", 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLParamFloat(r, "num", 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'num' is required"))
				Expect(value).To(Equal(float64(0)))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				ctx.URLParams.Add("num", "number")

				value, err := httpr.URLParamFloat(r, "num", 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'num' is not valid float number"))
				Expect(value).To(Equal(float64(0)))
			})

			It("returns the provided value", func() {
				value := float64(200.10)
				ctx.URLParams.Add("num", "number")
				Expect(httpr.URLParamFloatOrValue(r, "num", 64, value)).To(Equal(value))
			})
		})
	})

	Describe("URLParamTime", func() {
		It("parses the values successfully", func() {
			now := time.Now()
			ctx.URLParams.Add("from", now.Format(time.RFC3339Nano))

			value, err := httpr.URLParamTime(r, "from", time.RFC3339Nano)
			Expect(err).To(BeNil())
			Expect(value).To(BeTemporally("==", now))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLParamTime(r, "from", time.RFC3339Nano)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("Parameter 'from' is required"))
				Expect(value.IsZero()).To(BeTrue())
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				ctx.URLParams.Add("from", "time")

				value, err := httpr.URLParamTime(r, "from", time.RFC3339Nano)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.ErrorResponse)
				Expect(ok).To(BeTrue())

				rrErr, ok := (rErr.Err).(*httpr.HTTPError)
				Expect(ok).To(BeTrue())

				Expect(rrErr.Message).To(Equal("Parameter 'from' is not valid date time"))
				Expect(rrErr.Details).To(HaveLen(1))
				Expect(rrErr.Details[0]).To(Equal(fmt.Sprintf("Expected date time format '%s'", time.RFC3339Nano)))
				Expect(value.IsZero()).To(BeTrue())
			})

			It("returns the provided value", func() {
				now := time.Now()
				ctx.URLParams.Add("from", "time")
				Expect(httpr.URLParamTimeOrValue(r, "num", time.RFC3339Nano, now)).To(BeTemporally("==", now))
			})
		})
	})
})
