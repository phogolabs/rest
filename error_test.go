package rho_test

import (
	"fmt"
	"strings"

	"github.com/gosuri/uitable"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho"
)

var _ = Describe("Error", func() {
	It("returns the error message correctly", func() {
		err := rho.NewError(201, "Oh no!", "Unexpected error")
		Expect(err.Message).To(Equal("Oh no!"))
		Expect(err.Details).To(HaveLen(1))
		Expect(err.Details).To(ContainElement("Unexpected error"))
	})

	It("wraps the error successfully", func() {
		err := rho.NewError(201, "Oh no!", "Unexpected error")
		err.Wrap(fmt.Errorf("Inner Error"))
		Expect(err.Reason).To(MatchError("Inner Error"))
	})

	It("returns the correct error message", func() {

		err := rho.NewError(201, "Oh no!", "Unexpected error")
		err.Wrap(fmt.Errorf("Inner Error"))

		table := uitable.New()
		table.MaxColWidth = 80
		table.Wrap = true

		table.AddRow("code:", fmt.Sprintf("%d", err.Code))
		table.AddRow("message:", err.Message)
		table.AddRow("details:", strings.Join(err.Details, ", "))
		table.AddRow("reason:", err.Reason.Error())

		Expect(err.Error()).To(Equal(table.String()))
	})
})
