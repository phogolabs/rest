package httputil_test

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRho(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTPUtil Suite")
}

type T struct {
	Err  string `json:"err"`
	Name string `json:"name"`
}

func (t T) Bind(r *http.Request) error {
	if t.Err != "" {
		return fmt.Errorf(t.Err)
	}
	return nil
}
