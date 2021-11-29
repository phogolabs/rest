package rest_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/errors"
	"github.com/go-playground/validator/v10"
	"github.com/goware/errorx"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/phogolabs/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	var (
		recorder *httptest.ResponseRecorder
		request  *http.Request
		handle   handlerFn
		decode   decodeFn
	)

	ItHandlesTheError := func() {
		Context("when the error is unknown", func() {
			It("responds with error", func() {
				handle(recorder, request, fmt.Errorf("oh no!"))

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				err := &errorx.Errorx{}
				Expect(decode(err)).To(Succeed())

				Expect(err.Code).To(Equal(http.StatusInternalServerError))
				Expect(err.Message).To(Equal(http.StatusText(http.StatusInternalServerError)))
				Expect(err.Details).To(HaveLen(1))
				Expect(err.Details).To(ContainElement("oh no!"))
			})
		})

		Context("when is used in a middleware", func() {
			fn := func(next http.Handler) http.Handler {
				handle := func(w http.ResponseWriter, r *http.Request) {
					err := fmt.Errorf("oh no!")
					rest.Status(r, http.StatusUnauthorized)
					handle(w, r, err)
				}

				return http.HandlerFunc(handle)
			}

			handle := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			}

			It("returns the status code", func() {
				r := chi.NewMux()
				r.Use(fn)
				r.Post("/v1/users", handle)

				r.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusUnauthorized))

				err := &errorx.Errorx{}
				Expect(decode(err)).To(Succeed())

				Expect(err.Code).To(Equal(http.StatusUnauthorized))
				Expect(err.Message).To(Equal(http.StatusText(http.StatusUnauthorized)))
				Expect(err.Details).To(HaveLen(1))
				Expect(err.Details).To(ContainElement("oh no!"))
			})
		})

		Context("when the status is set", func() {
			It("responds with error", func() {
				rest.Status(request, http.StatusForbidden)
				handle(recorder, request, fmt.Errorf("oh no!"))

				Expect(recorder.Code).To(Equal(http.StatusForbidden))

				err := &errorx.Errorx{}
				Expect(decode(err)).To(Succeed())

				Expect(err.Code).To(Equal(http.StatusForbidden))
				Expect(err.Message).To(Equal(http.StatusText(http.StatusForbidden)))
				Expect(err.Details).To(HaveLen(1))
				Expect(err.Details).To(ContainElement("oh no!"))
			})
		})

		Context("when the error is chained", func() {
			It("responds with error", func() {
				rerr := errors.New("oh no!").AddTag("status", http.StatusRequestTimeout)
				handle(recorder, request, rerr)

				Expect(recorder.Code).To(Equal(http.StatusRequestTimeout))

				err := &errorx.Errorx{}
				Expect(decode(err)).To(Succeed())

				Expect(err.Code).To(Equal(http.StatusRequestTimeout))
				Expect(err.Message).To(Equal(http.StatusText(http.StatusRequestTimeout)))
				Expect(err.Details).To(HaveLen(1))
				Expect(err.Details).To(ContainElement("oh no!"))
			})
		})

		Context("when the error is multi error", func() {
			It("responds with error", func() {
				var errs error

				errs = multierror.Append(errs, fmt.Errorf("oh no!"))
				errs = multierror.Append(errs, fmt.Errorf("oh yes!"))

				handle(recorder, request, errs)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				err := &errorx.Errorx{}
				Expect(decode(err)).To(Succeed())

				Expect(err.Code).To(Equal(http.StatusInternalServerError))
				Expect(err.Message).To(Equal(http.StatusText(http.StatusInternalServerError)))
				Expect(err.Details).To(HaveLen(2))
				Expect(err.Details).To(ContainElement("oh no!"))
				Expect(err.Details).To(ContainElement("oh yes!"))
			})
		})

		Context("when the error is validation error", func() {
			It("responds with error", func() {
				entity := Person{}

				verr := validator.New().Struct(&entity)
				Expect(verr).NotTo(BeNil())

				handle(recorder, request, verr)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				err := &errorx.Errorx{}
				Expect(decode(err)).To(Succeed())

				Expect(err.Code).To(Equal(http.StatusInternalServerError))
				Expect(err.Message).To(Equal(http.StatusText(http.StatusInternalServerError)))
				Expect(err.Details).To(HaveLen(1))
				Expect(err.Details).To(ContainElement(verr.Error()))
			})
		})
	}

	Context("application/json", func() {
		BeforeEach(func() {
			request = NewJSONRequest(nil)
			recorder = httptest.NewRecorder()
			decode = json.NewDecoder(recorder.Body).Decode
			handle = rest.Error
		})

		ItHandlesTheError()
	})

	Context("application/json", func() {
		BeforeEach(func() {
			request = NewJSONRequest(nil)
			recorder = httptest.NewRecorder()
			decode = json.NewDecoder(recorder.Body).Decode
			handle = rest.ErrorJSON
		})

		ItHandlesTheError()
	})

	Context("application/xml", func() {
		BeforeEach(func() {
			request = NewXMLRequest(nil)
			recorder = httptest.NewRecorder()
			decode = xml.NewDecoder(recorder.Body).Decode
			handle = rest.ErrorXML
		})

		ItHandlesTheError()
	})
})
