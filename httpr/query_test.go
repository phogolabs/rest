package httpr_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/http/httpr"
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
			Expect(httpr.URLQueryParam(r, "name")).To(Equal("Jack"))
		})
	})

	Describe("URLQueryParamUUID", func() {
		It("parses the values successfully", func() {
			id := uuid.NewV4()
			addQueryParam("id", id.String())

			value, err := httpr.URLQueryParamUUID(r, "id")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(id))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLQueryParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'id' is required"))
				Expect(value).To(Equal(uuid.Nil))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("id", "wrong-uuid")

				value, err := httpr.URLQueryParamUUID(r, "id")
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'id' is not valid UUID"))
				Expect(value).To(Equal(uuid.Nil))
			})

			It("returns a nil value", func() {
				addQueryParam("id", "wrong-uuid")
				Expect(httpr.URLQueryParamUUIDOrNil(r, "id")).To(Equal(uuid.Nil))
			})

			It("returns the provided value", func() {
				id := uuid.NewV4()
				addQueryParam("id", "wrong-uuid")
				Expect(httpr.URLQueryParamUUIDOrValue(r, "id", id)).To(Equal(id))
			})
		})
	})

	Describe("URLQueryParamInt", func() {
		It("parses the values successfully", func() {
			num := int64(123)
			addQueryParam("num", "123")

			value, err := httpr.URLQueryParamInt(r, "num", 0, 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLQueryParamInt(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'num' is required"))
				Expect(value).To(Equal(int64(0)))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("num", "number")

				value, err := httpr.URLQueryParamInt(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'num' is not valid integer number"))
				Expect(value).To(Equal(int64(0)))
			})

			It("returns the provided value", func() {
				value := int64(200)
				addQueryParam("num", "number")
				Expect(httpr.URLQueryParamIntOrValue(r, "num", 0, 64, value)).To(Equal(value))
			})
		})
	})

	Describe("URLQueryParamUint", func() {
		It("parses the values successfully", func() {
			num := uint64(123)
			addQueryParam("num", "123")

			value, err := httpr.URLQueryParamUint(r, "num", 0, 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is negative", func() {
			It("parses the values successfully", func() {
				addQueryParam("num", "-123")

				value, err := httpr.URLQueryParamUint(r, "num", 0, 64)
				Expect(err).NotTo(BeNil())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'num' is not valid unsigned integer number"))
				Expect(value).To(Equal(uint64(0)))
			})
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLQueryParamUint(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'num' is required"))
				Expect(value).To(Equal(uint64(0)))
			})

			It("returns the provided value", func() {
				value := uint64(200)
				addQueryParam("num", "number")
				Expect(httpr.URLQueryParamUintOrValue(r, "num", 0, 64, value)).To(Equal(value))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("num", "number")

				value, err := httpr.URLQueryParamUint(r, "num", 0, 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'num' is not valid unsigned integer number"))
				Expect(value).To(Equal(uint64(0)))
			})
		})
	})

	Describe("URLQueryParamFloat", func() {
		It("parses the values successfully", func() {
			num := float64(123.11)
			addQueryParam("num", "123.11")

			value, err := httpr.URLQueryParamFloat(r, "num", 64)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(num))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLQueryParamFloat(r, "num", 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'num' is required"))
				Expect(value).To(Equal(float64(0)))
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("num", "number")

				value, err := httpr.URLQueryParamFloat(r, "num", 64)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'num' is not valid float number"))
				Expect(value).To(Equal(float64(0)))
			})

			It("returns the provided value", func() {
				value := float64(200.10)
				Expect(httpr.URLQueryParamFloatOrValue(r, "num", 64, value)).To(Equal(value))
			})
		})
	})

	Describe("URLQueryParamTime", func() {
		It("parses the values successfully", func() {
			now := time.Now()
			addQueryParam("from", now.Format(time.RFC3339Nano))

			value, err := httpr.URLQueryParamTime(r, "from", time.RFC3339Nano)
			Expect(err).To(BeNil())
			Expect(value).To(BeTemporally("==", now))
		})

		Context("when the parameter is missing", func() {
			It("returns an error response", func() {
				value, err := httpr.URLQueryParamTime(r, "from", time.RFC3339Nano)
				Expect(err).To(HaveOccurred())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'from' is required"))
				Expect(value.IsZero()).To(BeTrue())
			})
		})

		Context("when the parameter is malformed", func() {
			It("returns an error response", func() {
				addQueryParam("from", "time")

				value, err := httpr.URLQueryParamTime(r, "from", time.RFC3339Nano)
				Expect(err).To(HaveOccurred())
				Expect(value.IsZero()).To(BeTrue())

				rErr, ok := (err).(*httpr.Error)
				Expect(ok).To(BeTrue())

				Expect(rErr).To(MatchError("query parameter 'from' is not valid date time"))
				Expect(rErr.Details).To(HaveLen(1))
				Expect(rErr.Details[0]).To(Equal(fmt.Sprintf("expected date time format '%s'", time.RFC3339Nano)))
			})

			It("returns the provided value", func() {
				now := time.Now()
				addQueryParam("from", "time")
				Expect(httpr.URLQueryParamTimeOrValue(r, "num", time.RFC3339Nano, now)).To(BeTemporally("==", now))
			})
		})
	})
})
