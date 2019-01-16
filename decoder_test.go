package rest_test

import (
	"net/url"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-playground/errors"
	"github.com/phogolabs/rest"
	validator "gopkg.in/go-playground/validator.v9"
)

var _ = Describe("Decoder", func() {
	Describe("Content-Type: application/x-www-form-urlencoded", func() {
		It("decodes a form request successfully", func() {
			v := url.Values{}
			v.Add("name", "John")
			v.Add("age", "22")

			entity := Person{}
			request := NewFormRequest(v)

			Expect(rest.Decode(request, &entity)).To(Succeed())
			Expect(entity.Name).To(Equal("John"))
			Expect(entity.Age).To(BeEquivalentTo(22))
		})

		Context("when the parsing fails", func() {
			It("returns a n error", func() {
				v := url.Values{}
				v.Add("age", "-22")

				entity := Person{}
				request := NewFormRequest(v)

				err := errors.Cause(rest.Decode(request, &entity))
				Expect(err).To(MatchError("Field Namespace:age ERROR:Invalid Unsigned Integer Value '-22' Type 'uint' Namespace 'age'"))
			})
		})
	})

	Describe("Validate", func() {
		Context("when custom validation is registered", func() {
			BeforeEach(func() {
				rest.RegisterValidation("phone", func(field validator.FieldLevel) bool {
					phoneRegexp := regexp.MustCompile("\\+[0-9]+")
					value := field.Field().String()
					return phoneRegexp.MatchString(value)
				})
			})

			It("validates the entity with it", func() {
				v := url.Values{}
				v.Add("phone", "555ZERO")

				entity := Contact{}
				request := NewFormRequest(v)

				err := errors.Cause(rest.Decode(request, &entity))
				Expect(err).To(MatchError("Key: 'Contact.phone' Error:Field validation for 'phone' failed on the 'phone' tag"))
			})
		})

		Context("when the validation fails", func() {
			It("returns an error", func() {
				v := url.Values{}
				v.Add("age", "18")

				entity := Person{}
				request := NewFormRequest(v)

				err := errors.Cause(rest.Decode(request, &entity))
				Expect(err).To(MatchError("Key: 'Person.age' Error:Field validation for 'age' failed on the 'gte' tag"))
			})
		})
	})
})
