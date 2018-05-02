package httperr_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httperr"
)

var _ = Describe("Stack", func() {
	It("formats the stack correctly", func() {
		stack := httperr.NewStack()
		Expect(fmt.Sprintf("%+v", stack)).To(ContainSubstring("ginkgo"))
	})
})
