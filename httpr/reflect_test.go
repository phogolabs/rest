package httpr_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/http/httpr"
)

var _ = Describe("Stack", func() {
	It("formats the stack correctly", func() {
		stack := httpr.NewStack()
		Expect(fmt.Sprintf("%+v", stack)).To(ContainSubstring("ginkgo"))
	})
})
