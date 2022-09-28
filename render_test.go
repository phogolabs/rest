package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-playground/errors"
	"github.com/phogolabs/rest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bind", func() {
	var (
		response *Response
		request  *http.Request
	)

	BeforeEach(func() {
		response = &Response{}
		request = NewJSONRequest(nil)
	})

	It("binds a request successfully", func() {
		Expect(rest.Bind(request, response)).To(Succeed())
	})

	Context("when the binding fails", func() {
		BeforeEach(func() {
			response.Error = fmt.Errorf("oh no!")
		})

		It("returns an error", func() {
			err := rest.Bind(request, response)
			Expect(err).To(HaveOccurred())

			err = errors.Cause(err)
			Expect(err).To(MatchError("oh no!"))
		})
	})
})

var _ = Describe("Render", func() {
	var (
		response *Response
		request  *http.Request
	)

	BeforeEach(func() {
		response = &Response{}
		request = NewJSONRequest(nil)
	})

	It("renders a response successfully", func() {
		Expect(rest.Render(httptest.NewRecorder(), request, response)).To(Succeed())
	})

	Context("when the rendering fails", func() {
		BeforeEach(func() {
			response.Error = fmt.Errorf("oh no!")
		})

		It("returns an error", func() {
			err := rest.Render(httptest.NewRecorder(), request, response)
			Expect(err).To(HaveOccurred())

			err = errors.Cause(err)
			Expect(err).To(MatchError("oh no!"))
		})
	})
})

var _ = Describe("GetLogger", func() {
	It("returns a logger", func() {
		request := httptest.NewRequest("POST", "http://example.com", nil)
		Expect(rest.GetLogger(request)).NotTo(BeNil())
	})
})

var _ = Describe("Status", func() {
	It("sets the status", func() {
		request := httptest.NewRequest("POST", "http://example.com", nil)
		rest.Status(request, http.StatusOK)
	})
})
