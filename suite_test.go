package rest_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Person struct {
	Name string `form:"name" json:"name"`
	Age  uint   `form:"age" validate:"gte=21" json:"age"`
}

type Contact struct {
	Phone string `form:"phone" validate:"phone" json:"phone"`
}

func TestREST(t *testing.T) {
	log.SetOutput(GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest Suite")
}

func NewFormRequest(v url.Values) *http.Request {
	r := httptest.NewRequest("POST", "http://example.com", strings.NewReader(v.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func NewJSONRequest() *http.Request {
	r := httptest.NewRequest("POST", "http://example.com", nil)
	r.Header.Add("Content-Type", "application/josn")
	return r
}
