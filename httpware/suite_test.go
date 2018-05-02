package httpware_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttpware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpware Suite")
}
