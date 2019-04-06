package rest_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rest"
	validator "gopkg.in/go-playground/validator.v9"
)

type (
	handlerFn = func(http.ResponseWriter, *http.Request, error)
	decodeFn  = func(data interface{}) error
)

type Response struct {
	Error error
}

func (b *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return b.Error
}

func (b *Response) Bind(r *http.Request) error {
	return b.Error
}

type Person struct {
	Name    string `form:"name" xml:"name" json:"name,omitempty"`
	Age     uint   `form:"age" xml:"age" json:"age,omitemoty" validate:"gte=21" `
	Address string `default:"london"`
}

type Contact struct {
	Phone string `form:"phone" xml:"phone" json:"phone,omitempty" validate:"phone" `
}

func TestREST(t *testing.T) {
	log.SetOutput(GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest Suite")
}

var _ = BeforeSuite(func() {
	rest.RegisterValidation("phone", func(field validator.FieldLevel) bool {
		phoneRegexp := regexp.MustCompile("\\+[0-9]+")
		value := field.Field().String()
		return phoneRegexp.MatchString(value)
	})
})

func NewJSONRequest(data interface{}) *http.Request {
	buffer := &bytes.Buffer{}

	if data != nil {
		Expect(json.NewEncoder(buffer).Encode(data)).To(Succeed())
	}

	r := httptest.NewRequest("POST", "http://example.com", buffer)
	r.Header.Add("Content-Type", "application/json")
	return r
}

func NewXMLRequest(data interface{}) *http.Request {
	buffer := &bytes.Buffer{}

	if data != nil {
		Expect(xml.NewEncoder(buffer).Encode(data)).To(Succeed())
	}

	r := httptest.NewRequest("POST", "http://example.com", buffer)
	r.Header.Add("Content-Type", "application/xml")
	return r
}

func NewGobRequest(data interface{}) *http.Request {
	buffer := &bytes.Buffer{}

	if data != nil {
		Expect(gob.NewEncoder(buffer).Encode(data)).To(Succeed())
	}

	r := httptest.NewRequest("POST", "http://example.com", buffer)
	r.Header.Add("Content-Type", "application/gob")
	return r
}

func NewFormRequest(v url.Values) *http.Request {
	r := httptest.NewRequest("POST", "/v1/users", strings.NewReader(v.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return r
}
