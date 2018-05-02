package httperr_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/apex/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httperr"
)

var _ = Describe("Error", func() {
	It("returns the error message correctly", func() {
		err := httperr.New(201, "Oh no!", "Unexpected error")
		Expect(err.Message).To(Equal("Oh no!"))
		Expect(err.Details).To(HaveLen(1))
		Expect(err.Details).To(ContainElement("Unexpected error"))
	})

	It("reports the stack trace correctly", func() {
		err := httperr.New(201, "Oh no!", "Unexpected error")
		Expect(fmt.Sprintf("%v", err.StackTrace())).To(ContainSubstring("error_test.go"))
	})

	It("returns the fields correctly", func() {
		leaf := fmt.Errorf("Inner Error")

		err := httperr.New(201, "Oh no!", "Unexpected error")
		err.Wrap(httperr.New(202, "Madness").Wrap(leaf))

		fields := err.Fields()

		Expect(fields).To(HaveKeyWithValue("code", 201))
		Expect(fields).NotTo(HaveKey("message"))
		Expect(fields).To(HaveKeyWithValue("details[0]", "Unexpected error"))
		Expect(fields).To(HaveKey("reason"))

		rfields, ok := fields["reason"].(log.Fields)
		Expect(ok).To(BeTrue())

		Expect(rfields).To(HaveKeyWithValue("code", 202))
		Expect(rfields).To(HaveKeyWithValue("message", "Madness"))
		Expect(rfields).NotTo(HaveKey("details[0]"))
		Expect(rfields).To(HaveKeyWithValue("reason", "Inner Error"))
	})

	It("wraps the error successfully", func() {
		err := httperr.New(201, "Oh no!", "Unexpected error")
		err.Wrap(fmt.Errorf("Inner Error"))
		Expect(err.Reason).To(MatchError("Inner Error"))
	})

	It("returns the error that causes the error", func() {
		leaf := fmt.Errorf("Inner Error")

		err := httperr.New(201, "Oh no!", "Unexpected error")
		err.Wrap(httperr.New(202, "Madness").Wrap(leaf))

		Expect(err.Cause()).To(MatchError("Inner Error"))
	})

	Context("when there is not reason", func() {
		It("returns the error that causes the error", func() {
			err := httperr.New(201, "Oh no!", "Unexpected error")
			Expect(err.Cause()).To(Equal(err))
		})
	})

	It("returns the correct error message", func() {
		err := httperr.New(201, "Oh no!", "Unexpected error")
		err.Wrap(fmt.Errorf("Inner Error"))
		Expect(err.Error()).To(Equal("Oh no!"))
	})
})

var _ = Describe("HandleErr", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the error is regular error", func() {
		It("handles the error", func() {
			err := fmt.Errorf("Oh no!")
			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httperr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Internal Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the error is ErrorReponse", func() {
		It("handles the error", func() {
			err := &httperr.Response{
				StatusCode: http.StatusBadGateway,
				Err:        httperr.New(1, "Oh no!"),
			}

			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(err.StatusCode))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(1)))
			Expect(payload).To(HaveKeyWithValue("message", "Oh no!"))
			Expect(payload).NotTo(HaveKey("reason"))
		})
	})
})
