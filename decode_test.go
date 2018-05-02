package rho_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho"
)

var _ = Describe("Decode", func() {
	var (
		r    *http.Request
		body *bytes.Buffer
	)

	BeforeEach(func() {
		body = &bytes.Buffer{}
		r = httptest.NewRequest("GET", "http://example.com", body)
		r.Header.Add("Content-Type", "application/json")
	})

	It("descodes the request successfully", func() {
		t := T{Name: "Jack"}
		Expect(json.NewEncoder(body).Encode(&t)).To(Succeed())

		t2 := T{}
		Expect(rho.Decode(r, &t2)).To(Succeed())
		Expect(t2).To(Equal(t))
	})

	Context("when the decoder fails", func() {
		It("returns the error", func() {
			body.WriteString("wrong")

			t := time.Now()
			Expect(rho.Decode(r, &t)).To(MatchError("Unable to unmarshal request body"))
		})
	})

	Context("when the binder fails", func() {
		It("returns the error", func() {
			t := T{Name: "Jack", Err: "Oh no"}
			Expect(json.NewEncoder(body).Encode(&t)).To(Succeed())

			t2 := T{}
			Expect(rho.Decode(r, &t2)).To(MatchError("Unable to bind request"))
		})
	})

	Context("when the validation fails", func() {
		It("returns the error", func() {
			t := time.Now()
			Expect(json.NewEncoder(body).Encode(&t)).To(Succeed())

			var t2 time.Time
			Expect(rho.Decode(r, &t2)).To(MatchError("Unable to validate request"))
		})
	})
})
