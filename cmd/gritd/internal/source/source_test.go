package source_test

import (
	. "github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/cmd/gritd/internal/source/githubdriver"
	"github.com/gritcli/grit/internal/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func New()", func() {
	DescribeTable(
		"it returns a source based on the driver and configuration",
		func(cfg config.Source, expectType Source) {
			src, err := New(cfg)
			Expect(err).ShouldNot(HaveOccurred())
			defer src.Close()

			Expect(src).To(BeAssignableToTypeOf(expectType))
		},
		Entry(
			"GitHub driver",
			config.Source{
				Name: "github.com",
				Config: config.GitHubConfig{
					Domain: "github.com",
				},
			},
			&githubdriver.Source{},
		),
	)
})
