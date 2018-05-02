package rho_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho"
	"github.com/phogolabs/rho/httperr"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Query", func() {
	var r *http.Request

	BeforeEach(func() {
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	addQueryParam := func(key, value string) {
		q := r.URL.Query()
		q.Add(key, value)
		r.URL.RawQuery = q.Encode()
	}

	Describe("URLQueryParam", func() {
		It("returns the value successfully", func() {
			addQueryParam("name", "Jack")
			Expect(rho.URLQueryParam(r, "name")).To(Equal("Jack"))
		})
	})

	Describe("URLQueryParamUUID", func() {
		It("parses the values successfully", func() {
			id := uuid.NewV4()
			addQueryParam("id", id.String())

			value, err := rho.URLQueryParamUUID(r, "id")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(id))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := rho.URLQueryParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'id' is required"))
				Expect(value).To(Equal(uuid.Nil))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("id", "wrong-uuid")

				value, err := rho.URLQueryParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'id' is not valid UUID"))
				Expect(value).To(Equal(uuid.Nil))
			})

			It("returns a nil value", func() {
				addQueryParam("id", "wrong-uuid")
				Expect(rho.URLQueryParamUUIDOrNil(r, "id")).To(Equal(uuid.Nil))
			})

			It("returns the provided value", func() {
				id := uuid.NewV4()
				addQueryParam("id", "wrong-uuid")
				Expect(rho.URLQueryParamUUIDOrValue(r, "id", id)).To(Equal(id))
			})
		})
	})

	Describe("URLQueryParamInt", func() {
		It("parses the values successfully", func() {
			num := int64(123)
			addQueryParam("num", "123")

			value, err := rho.URLQueryParamInt(r, "num", 0, 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := rho.URLQueryParamInt(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'num' is required"))
				Expect(value).To(Equal(int64(0)))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("num", "number")

				value, err := rho.URLQueryParamInt(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'num' is not valid integer number"))
				Expect(value).To(Equal(int64(0)))
			})

			It("returns the provided value", func() {
				value := int64(200)
				addQueryParam("num", "number")
				Expect(rho.URLQueryParamIntOrValue(r, "num", 0, 64, value)).To(Equal(value))
			})
		})
	})

	Describe("URLQueryParamUint", func() {
		It("parses the values successfully", func() {
			num := uint64(123)
			addQueryParam("num", "123")

			value, err := rho.URLQueryParamUint(r, "num", 0, 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is negative", func() {
			It("parses the values successfully", func() {
				addQueryParam("num", "-123")

				value, err := rho.URLQueryParamUint(r, "num", 0, 64)
				Expect(err).NotTo(BeNil())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'num' is not valid unsigned integer number"))
				Expect(value).To(Equal(uint64(0)))
			})
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := rho.URLQueryParamUint(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'num' is required"))
				Expect(value).To(Equal(uint64(0)))
			})

			It("returns the provided value", func() {
				value := uint64(200)
				addQueryParam("num", "number")
				Expect(rho.URLQueryParamUintOrValue(r, "num", 0, 64, value)).To(Equal(value))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("num", "number")

				value, err := rho.URLQueryParamUint(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'num' is not valid unsigned integer number"))
				Expect(value).To(Equal(uint64(0)))
			})
		})
	})

	Describe("URLQueryParamFloat", func() {
		It("parses the values successfully", func() {
			num := float64(123.11)
			addQueryParam("num", "123.11")

			value, err := rho.URLQueryParamFloat(r, "num", 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := rho.URLQueryParamFloat(r, "num", 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'num' is required"))
				Expect(value).To(Equal(float64(0)))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("num", "number")

				value, err := rho.URLQueryParamFloat(r, "num", 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'num' is not valid float number"))
				Expect(value).To(Equal(float64(0)))
			})

			It("returns the provided value", func() {
				value := float64(200.10)
				// addQueryParam("num", "number")
				Expect(rho.URLQueryParamFloatOrValue(r, "num", 64, value)).To(Equal(value))
			})
		})
	})

	Describe("URLQueryParamTime", func() {
		It("parses the values successfully", func() {
			now := time.Now()
			addQueryParam("from", now.Format(time.RFC3339Nano))

			value, err := rho.URLQueryParamTime(r, "from", time.RFC3339Nano)
			Expect(err).To(BeNil())
			Expect(value).To(BeTemporally("==", now))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := rho.URLQueryParamTime(r, "from", time.RFC3339Nano)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'from' is required"))
				Expect(value.IsZero()).To(BeTrue())
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("from", "time")

				value, err := rho.URLQueryParamTime(r, "from", time.RFC3339Nano)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httperr.Response)
				Expect(ok).To(BeTrue())

				Expect(rErr.Err.Message).To(Equal("Query Parameter 'from' is not valid date time"))
				Expect(rErr.Err.Details).To(HaveLen(1))
				Expect(rErr.Err.Details[0]).To(Equal(fmt.Sprintf("Expected date time format '%s'", time.RFC3339Nano)))
				Expect(value.IsZero()).To(BeTrue())
			})

			It("returns the provided value", func() {
				now := time.Now()
				addQueryParam("from", "time")
				Expect(rho.URLQueryParamTimeOrValue(r, "num", time.RFC3339Nano, now)).To(BeTemporally("==", now))
			})
		})
	})
})
