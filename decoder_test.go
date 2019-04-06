package rest_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-playground/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rest"
)

var _ = Describe("Decode", func() {
	var request *http.Request

	Describe("JSON", func() {
		BeforeEach(func() {
			contact := &Contact{Phone: "+188123451"}
			request = NewJSONRequest(contact)
		})

		It("decodes a form request successfully", func() {
			entity := Contact{}

			Expect(rest.Decode(request, &entity)).To(Succeed())
			Expect(entity.Phone).To(Equal("+188123451"))
		})
	})

	Describe("XML", func() {
		var contact *Contact

		BeforeEach(func() {
			contact = &Contact{Phone: "+188123451"}
		})

		JustBeforeEach(func() {
			request = NewXMLRequest(contact)
		})

		It("decodes a form request successfully", func() {
			entity := Contact{}

			Expect(rest.Decode(request, &entity)).To(Succeed())
			Expect(entity.Phone).To(Equal("+188123451"))
		})

		Context("when the validation fails", func() {
			BeforeEach(func() {
				contact = &Contact{Phone: "088HIPPO"}
			})

			It("returns an error", func() {
				entity := Contact{}

				err := rest.Decode(request, &entity)
				Expect(err).To(HaveOccurred())

				err = errors.Cause(err)
				Expect(err).To(MatchError("Key: 'Contact.phone' Error:Field validation for 'phone' failed on the 'phone' tag"))
			})
		})
	})

	Describe("FORM", func() {
		BeforeEach(func() {
			v := url.Values{}
			v.Add("name", "John")
			v.Add("age", "22")

			request = NewFormRequest(v)
		})

		It("decodes a form request successfully", func() {
			entity := Person{}

			Expect(rest.Decode(request, &entity)).To(Succeed())
			Expect(entity.Name).To(Equal("John"))
			Expect(entity.Age).To(BeEquivalentTo(22))
		})
	})

	Context("when the Content-Tyoe is UNKNOWN", func() {
		BeforeEach(func() {
			contact := &Contact{Phone: "088HIPPO"}
			request = NewGobRequest(contact)
		})

		It("returns an error", func() {
			entity := Contact{}

			err := rest.Decode(request, &entity)
			Expect(err).To(HaveOccurred())

			err = errors.Cause(err)
			Expect(err).To(MatchError("render: unable to automatically decode the request content type"))
		})
	})
})

var _ = Describe("DecodePath", func() {
	var request *http.Request

	type User struct {
		ID int `path:"id"`
	}

	BeforeEach(func() {
		request = httptest.NewRequest("POST", "http://example.com/users/1", nil)
	})

	It("decodes the request successfully", func() {
		router := chi.NewMux()

		router.Mount("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := User{}
			Expect(rest.DecodePath(r, &user)).To(Succeed())
			Expect(user.ID).To(Equal(1))
		}))

		router.ServeHTTP(httptest.NewRecorder(), request)
	})

	Context("when the types are incompatible", func() {
		BeforeEach(func() {
			request = httptest.NewRequest("POST", "http://example.com/users/root", nil)
		})

		It("returns an error", func() {
			router := chi.NewMux()

			router.Mount("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := User{}
				Expect(rest.DecodePath(r, &user)).To(MatchError("Field Namespace:id ERROR:Invalid Integer Value 'root' Type 'int' Namespace 'id'"))
			}))

			router.ServeHTTP(httptest.NewRecorder(), request)
		})
	})
})

var _ = Describe("DecodeQuery", func() {
	var request *http.Request

	type User struct {
		ID int `query:"id"`
	}

	BeforeEach(func() {
		request = httptest.NewRequest("POST", "http://example.com/?id=1", nil)
	})

	It("decodes the request successfully", func() {
		user := User{}
		Expect(rest.DecodeQuery(request, &user)).To(Succeed())
		Expect(user.ID).To(Equal(1))
	})

	Context("when the types are incompatible", func() {
		BeforeEach(func() {
			request = httptest.NewRequest("POST", "http://example.com/?id=root", nil)
		})

		It("returns an error", func() {
			user := User{}
			Expect(rest.DecodeQuery(request, &user)).To(MatchError("Field Namespace:id ERROR:Invalid Integer Value 'root' Type 'int' Namespace 'id'"))
		})
	})
})

var _ = Describe("DecodeHeader", func() {
	var request *http.Request

	type User struct {
		ID int `header:"X-User-Id"`
	}

	BeforeEach(func() {
		request = httptest.NewRequest("POST", "http://example.com", nil)
		request.Header.Set("X-User-Id", "1")
	})

	It("decodes the request successfully", func() {
		user := User{}
		Expect(rest.DecodeHeader(request, &user)).To(Succeed())
		Expect(user.ID).To(Equal(1))
	})

	Context("when the types are incompatible", func() {
		BeforeEach(func() {
			request.Header.Set("X-User-Id", "root")
		})

		It("returns an error", func() {
			user := User{}
			Expect(rest.DecodeHeader(request, &user)).To(MatchError("Field Namespace:X-User-Id ERROR:Invalid Integer Value 'root' Type 'int' Namespace 'X-User-Id'"))
		})
	})
})
