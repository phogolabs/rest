package rest_test

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rest"
)

var _ = Describe("Respond", func() {
	It("renders json", func() {
		request := NewJSONRequest(nil)
		recorder := httptest.NewRecorder()

		rest.Respond(recorder, request, "hello")

		Expect(recorder.Body.String()).To(ContainSubstring("hello"))
		Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json; charset=utf-8"))
	})
})

var _ = Describe("JSON", func() {
	It("renders json", func() {
		request := NewJSONRequest(nil)
		recorder := httptest.NewRecorder()

		rest.JSON(recorder, request, "hello")

		Expect(recorder.Body.String()).To(ContainSubstring("hello"))
		Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json; charset=utf-8"))
	})
})

var _ = Describe("XML", func() {
	It("renders json", func() {
		request := NewXMLRequest(nil)
		recorder := httptest.NewRecorder()

		rest.XML(recorder, request, "hello")

		Expect(recorder.Body.String()).To(ContainSubstring("hello"))
		Expect(recorder.Header().Get("Content-Type")).To(Equal("application/xml; charset=utf-8"))
	})
})

var _ = Describe("PlainText", func() {
	It("renders text", func() {
		request := NewJSONRequest(nil)
		recorder := httptest.NewRecorder()

		rest.PlainText(recorder, request, "hello")

		Expect(recorder.Body.String()).To(Equal("hello"))
		Expect(recorder.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
	})
})

var _ = Describe("Data", func() {
	It("renders data", func() {
		request := NewJSONRequest(nil)
		recorder := httptest.NewRecorder()

		rest.Data(recorder, request, []byte("hello"))

		Expect(recorder.Body.String()).To(Equal("hello"))
		Expect(recorder.Header().Get("Content-Type")).To(Equal("application/octet-stream"))
	})
})

var _ = Describe("HTML", func() {
	It("renders html", func() {
		request := NewJSONRequest(nil)
		recorder := httptest.NewRecorder()

		rest.HTML(recorder, request, "hello")

		Expect(recorder.Body.String()).To(Equal("hello"))
		Expect(recorder.Header().Get("Content-Type")).To(Equal("text/html; charset=utf-8"))
	})
})

var _ = Describe("NoContent", func() {
	It("sets the no content status code", func() {
		request := NewJSONRequest(nil)
		recorder := httptest.NewRecorder()

		rest.NoContent(recorder, request)
		Expect(recorder.Code).To(Equal(204))
	})
})

var _ = Describe("EncodeHeader", func() {
	type User struct {
		ID int `header:"X-User-Id"`
	}

	It("encodes the header successfully", func() {
		u := &User{ID: 2}
		r := httptest.NewRecorder()

		Expect(rest.EncodeHeader(r, u)).To(Succeed())
		Expect(r.Header().Get("X-User-Id")).To(Equal("2"))
	})
})
