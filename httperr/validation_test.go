package httperr_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httperr"
	validator "gopkg.in/go-playground/validator.v9"
)

var _ = Describe("Validation", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the error is validation error", func() {
		It("respond with the provided information", func() {
			type User struct {
				Name     string `json:"name" validate:"required"`
				Password string `json:"password" validate:"required"`
			}

			v := validator.New()
			err := v.Struct(&User{})

			httperr.Respond(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			Expect(w.Body.String()).To(ContainSubstring("Field 'Name' is not valid"))
			Expect(w.Body.String()).To(ContainSubstring("Field 'Password' is not valid"))
		})
	})

	Context("when the error is InvalidValidationError", func() {
		It("handles the error correctoy", func() {
			err := &validator.InvalidValidationError{Type: reflect.TypeOf(*r)}
			resp := httperr.ValidationError(err)
			Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			Expect(resp.Err.Message).To(Equal("Validation Error"))
		})
	})
})
