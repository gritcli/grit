package githubsource_test

import (
	. "github.com/gritcli/grit/daemon/internal/builtins/githubsource"
	"github.com/gritcli/grit/daemon/internal/builtins/gitvcs"
	"github.com/gritcli/grit/daemon/internal/driver/configtest"
	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
	. "github.com/onsi/ginkgo/v2"
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

var _ = Describe("type configLoader", func() {
	configtest.TestSourceDriver(
		Registration,
		Config{},
		[]vcsdriver.Registration{
			gitvcs.Registration,
		},
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
