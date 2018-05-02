package httperr_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httperr"
)

var _ = Describe("Util", func() {
	It("returns the ParamRequired response", func() {
		err := httperr.ParamRequired("id")
		Expect(err).NotTo(BeNil())

		resp, ok := err.(*httperr.Response)
		Expect(ok).To(BeTrue())
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(resp.Err.Code).To(Equal(httperr.CodeParamRequired))
		Expect(resp.Err.Message).To(Equal("Parameter 'id' is required"))
	})

	It("returns the ParamParse response", func() {
		err := httperr.ParamParse("id", "string", fmt.Errorf("Oh no!"))
		Expect(err).NotTo(BeNil())

		resp, ok := err.(*httperr.Response)
		Expect(ok).To(BeTrue())
		Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		Expect(resp.Err.Code).To(Equal(httperr.CodeParamInvalid))
		Expect(resp.Err.Message).To(Equal("Parameter 'id' is not valid string"))
	})

	It("returns the QueryParamRequired response", func() {
		err := httperr.QueryParamRequired("id")
		Expect(err).NotTo(BeNil())

		resp, ok := err.(*httperr.Response)
		Expect(ok).To(BeTrue())
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(resp.Err.Code).To(Equal(httperr.CodeQueryParamRequired))
		Expect(resp.Err.Message).To(Equal("Query Parameter 'id' is required"))
	})

	It("returns the QueryParamParse response", func() {
		err := httperr.QueryParamParse("id", "string", fmt.Errorf("Oh no!"))
		Expect(err).NotTo(BeNil())

		resp, ok := err.(*httperr.Response)
		Expect(ok).To(BeTrue())
		Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		Expect(resp.Err.Code).To(Equal(httperr.CodeQueryParamInvalid))
		Expect(resp.Err.Message).To(Equal("Query Parameter 'id' is not valid string"))
	})
})
