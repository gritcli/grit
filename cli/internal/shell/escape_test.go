package shell_test

import (
	. "github.com/gritcli/grit/cli/internal/shell"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Escape()", func() {
	DescribeTable(
		"it escapes the string",
		func(input, expect string) {
			Expect(Escape(input)).To(Equal(expect))
		},
		Entry("empty", ``, `''`),
		Entry("plain string", `foo`, `'foo'`),
		Entry("string containing single quote", `foo'bar`, `'foo'"'"'bar'`),
	)
})
