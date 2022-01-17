package githubsource_test

import (
	"github.com/gritcli/grit/driver/configtest"
	. "github.com/gritcli/grit/driver/sourcedriver/githubsource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Config", func() {
	Describe("func DescribeSourceConfig()", func() {
		DescribeTable(
			"it describes the source",
			func(cfg Config, expect string) {
				Expect(cfg.DescribeSourceConfig()).To(Equal(expect))
			},
			Entry(
				"github.com",
				Config{Domain: "github.com"},
				"github.com",
			),
			Entry(
				"github enterprise server",
				Config{Domain: "code.example.com"},
				"code.example.com (github enterprise server)",
			),
		)
	})
})

var _ = Describe("type configSchema", func() {
	configtest.TestSourceDriver(
		Registration,
		Config{},
		configtest.SourceSuccess(
			"authentication token",
			`source "github" "github" {
				token = "<token>"
			}`,
			Config{
				Domain: "github.com",
				Token:  "<token>",
			},
		),
		configtest.SourceSuccess(
			"github enterprise server",
			`source "github" "github" {
				domain = "github.example.com"
			}`,
			Config{
				Domain: "github.example.com",
			},
		),
	)
})
