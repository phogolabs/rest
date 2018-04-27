package rho_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRho(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RHO Suite")
}
