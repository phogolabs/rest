package httpr_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/go-chi/chi/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger: log.New(GinkgoWriter, "", log.LstdFlags),
		})
}

func TestHTTPR(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTPR Suite")
}

type T struct {
	Err  string `json:"err"`
	Name string `json:"name" validate:"required"`
}

func (t *T) Bind(r *http.Request) error {
	if t.Err != "" {
		return fmt.Errorf(t.Err)
	}
	return nil
}

func unmarshalErrResponse(body *bytes.Buffer) map[string]interface{} {
	payload := make(map[string]interface{})
	Expect(json.NewDecoder(body).Decode(&payload)).To(Succeed())
	Expect(payload).To(HaveKey("error"))

	err, ok := payload["error"].(map[string]interface{})
	Expect(ok).To(BeTrue())
	return err
}
