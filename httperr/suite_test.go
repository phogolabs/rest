package httperr_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttperr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTPErr Suite")
}

func unmarshalErrResponse(body *bytes.Buffer) map[string]interface{} {
	payload := make(map[string]interface{})
	Expect(json.NewDecoder(body).Decode(&payload)).To(Succeed())
	Expect(payload).To(HaveKey("error"))

	err, ok := payload["error"].(map[string]interface{})
	Expect(ok).To(BeTrue())
	return err
}
