package httpr_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/http/httpr"
)

var _ = Describe("Error", func() {
	It("returns the error message correctly", func() {
		err := httpr.NewError(201, "Oh no!", "Unexpected error")
		Expect(err.Message).To(Equal("Oh no!"))
		Expect(err.Details).To(HaveLen(1))
		Expect(err.Details).To(ContainElement("Unexpected error"))
	})

	Context("when the code is zero", func() {
		It("sets internal code", func() {
			err := httpr.NewError(0, "Oh no!", "Unexpected error")
			Expect(err.Code).To(Equal(httpr.CodeInternal))
		})
	})

	It("reports the stack trace correctly", func() {
		err := httpr.NewError(201, "Oh no!", "Unexpected error")
		Expect(fmt.Sprintf("%v", err.StackTrace())).To(ContainSubstring("error_test.go"))
	})

	It("returns the fields correctly", func() {
		leaf := fmt.Errorf("Inner Error")

		err := httpr.NewError(201, "Oh no!", "Unexpected error")
		err.Wrap(httpr.NewError(202, "Madness").Wrap(leaf))

		fields := err.Fields()

		Expect(fields).To(HaveKeyWithValue("code", 201))
		Expect(fields).NotTo(HaveKey("message"))
		Expect(fields).To(HaveKeyWithValue("details[0]", "Unexpected error"))
		Expect(fields).To(HaveKey("reason"))

		rfields, ok := fields["reason"].(httpr.FieldsFormatter)
		Expect(ok).To(BeTrue())

		Expect(rfields).To(HaveKeyWithValue("code", 202))
		Expect(rfields).To(HaveKeyWithValue("message", "Madness"))
		Expect(rfields).NotTo(HaveKey("details[0]"))
		Expect(rfields).To(HaveKeyWithValue("reason", "Inner Error"))
	})

	It("wraps the error successfully", func() {
		err := httpr.NewError(201, "Oh no!", "Unexpected error")
		err.Wrap(fmt.Errorf("Inner Error"))
		Expect(err.Reason).To(MatchError("Inner Error"))
	})

	It("returns the error that causes the error", func() {
		leaf := fmt.Errorf("Inner Error")

		err := httpr.NewError(201, "Oh no!", "Unexpected error")
		err.Wrap(httpr.NewError(202, "Madness").Wrap(leaf))

		Expect(err.Cause()).To(MatchError("Inner Error"))
	})

	Context("when there is not reason", func() {
		It("returns the error that causes the error", func() {
			err := httpr.NewError(201, "Oh no!", "Unexpected error")
			Expect(err.Cause()).To(Equal(err))
		})
	})

	It("returns the correct error message", func() {
		err := httpr.NewError(201, "Oh no!", "Unexpected error")
		err.Wrap(fmt.Errorf("Inner Error"))
		Expect(err.Error()).To(Equal("Oh no!"))
	})
})

var _ = Describe("ErrorList", func() {
	It("returns the all error messages' fields", func() {
		m := httpr.ErrorList{}
		m = append(m, httpr.NewError(1, "Oh no!"))
		m = append(m, httpr.NewError(2, "Oh yes!"))

		err := httpr.NewError(1, "A lot of errors!")
		err.Wrap(m)

		f := err.Fields()
		Expect(f).To(HaveKeyWithValue("code", 1))
		Expect(f).To(HaveKey("reason"))
		Expect(f["reason"]).To(HaveKey("errors[0]"))
		Expect(f["reason"]).To(HaveKey("errors[1]"))
		Expect(f).NotTo(HaveKey("message"))
	})

	It("returns the all error messages", func() {
		m := httpr.ErrorList{}
		m = append(m, httpr.NewError(1, "Oh no!"))
		m = append(m, httpr.NewError(2, "Oh yes!"))
		Expect(m).To(MatchError("Oh no!;Oh yes!"))

		f := m.Fields()
		Expect(f).To(HaveKey("errors[0]"))
		Expect(f).To(HaveKey("errors[1]"))
	})
})

var _ = Describe("FieldsFormatter", func() {
	It("formats the fields successfully", func() {
		f := httpr.FieldsFormatter{"id": 1, "name": "root"}
		Expect(f.String()).To(Equal("[id:1 name:root]"))
	})

	It("adds a new field successfully", func() {
		f := httpr.FieldsFormatter{"id": 1, "name": "root"}
		f.Add("pass", "swordfish")
		Expect(f).To(HaveKeyWithValue("pass", "swordfish"))
	})
})

var _ = Describe("RenderError", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		formatter := &middleware.DefaultLogFormatter{
			Logger: log.New(GinkgoWriter, "", log.LstdFlags),
		}

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
		r = middleware.WithLogEntry(r, formatter.NewLogEntry(r))
	})

	Context("when the error is http error", func() {
		It("handles the error", func() {
			err := httpr.NewError(1, "Oh no!")
			httpr.RenderError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(1)))
			Expect(payload).To(HaveKeyWithValue("message", "Oh no!"))
		})

		Context("when the error is nested", func() {
			It("handles the error", func() {
				err := httpr.NewError(1, "Oh no!")
				err.Wrap(httpr.NewError(1, "Oh no!"))
				httpr.RenderError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(1)))
				Expect(payload).To(HaveKeyWithValue("message", "Oh no!"))
			})
		})
	})

	Context("when the error is nil", func() {
		It("handles the error", func() {
			httpr.RenderError(w, r, nil)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.Len()).To(BeZero())
		})
	})

	Context("when the error is regular error", func() {
		It("handles the error", func() {
			err := fmt.Errorf("Oh no!")
			httpr.RenderError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Internal Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})
