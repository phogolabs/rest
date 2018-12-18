package rest_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/go-chi/render"
	"github.com/go-playground/form"
	"github.com/goware/errorx"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/phogolabs/rest"
	validator "gopkg.in/go-playground/validator.v9"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StatusErr", func() {
	var r *http.Request

	BeforeEach(func() {
		r = NewFormRequest(url.Values{})
	})

	Context("when the err is a validation error", func() {
		It("sets status Unproccessable Entity: 422", func() {
			rest.StatusErr(r, validator.ValidationErrors{})

			code, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(code).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	Context("when the err is a form decoder error", func() {
		It("sets status Bad Request: 400", func() {
			rest.StatusErr(r, form.DecodeErrors{})

			code, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the err is a json field error", func() {
		It("sets status Bad Request: 400", func() {
			rest.StatusErr(r, &json.UnmarshalFieldError{})

			code, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the err is a json type error", func() {
		It("sets status Bad Request: 400", func() {
			rest.StatusErr(r, &json.UnmarshalTypeError{})

			code, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the err is a sql no rows error", func() {
		It("sets status Not Found: 404", func() {
			rest.StatusErr(r, sql.ErrNoRows)

			code, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(code).To(Equal(http.StatusNotFound))
		})
	})

	Context("when the err is unknown", func() {
		It("sets status Internal Server Errror: 500", func() {
			rest.StatusErr(r, fmt.Errorf("oh no!"))

			code, ok := r.Context().Value(render.StatusCtxKey).(int)
			Expect(ok).To(BeTrue())
			Expect(code).To(Equal(http.StatusInternalServerError))
		})
	})
})

var _ = Describe("Respond", func() {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		request = NewJSONRequest()
		recorder = httptest.NewRecorder()
	})

	Context("when the error is unknown", func() {
		It("responds with error", func() {
			rest.Respond(recorder, request, fmt.Errorf("oh no!"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			err := &errorx.Errorx{}
			Expect(json.NewDecoder(recorder.Body).Decode(err)).To(Succeed())

			Expect(err.Code).To(Equal(http.StatusInternalServerError))
			Expect(err.Message).To(Equal(http.StatusText(http.StatusInternalServerError)))
			Expect(err.Details).To(HaveLen(1))
			Expect(err.Details).To(ContainElement("oh no!"))
		})
	})

	Context("when the error is multi error", func() {
		It("responds with error", func() {
			var errs error

			errs = multierror.Append(errs, fmt.Errorf("oh no!"))
			errs = multierror.Append(errs, fmt.Errorf("oh yes!"))

			rest.Respond(recorder, request, errs)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			err := &errorx.Errorx{}
			Expect(json.NewDecoder(recorder.Body).Decode(err)).To(Succeed())

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

			rest.Respond(recorder, request, verr)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			err := &errorx.Errorx{}
			Expect(json.NewDecoder(recorder.Body).Decode(err)).To(Succeed())

			Expect(err.Code).To(Equal(http.StatusInternalServerError))
			Expect(err.Message).To(Equal(http.StatusText(http.StatusInternalServerError)))
			Expect(err.Details).To(HaveLen(1))
			Expect(err.Details).To(ContainElement(verr.Error()))
		})
	})
})
