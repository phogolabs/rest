package rho_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho"
)

var _ = Describe("Conv Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when a time error occurs", func() {
		It("handles the error correctly", func() {
			err := &time.ParseError{}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to parse date time"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when a num error occurs", func() {
		It("handles the error correctly", func() {
			err := &strconv.NumError{Err: fmt.Errorf("oh no!")}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to parse number"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})
