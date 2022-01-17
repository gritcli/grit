package gitvcs_test

import (
	. "github.com/gritcli/grit/driver/vcsdriver/gitvcs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Config", func() {
	Describe("func DescribeVCSConfig()", func() {
		DescribeTable(
			"it describes the source",
			func(cfg Config, expect string) {
				Expect(cfg.DescribeVCSConfig()).To(Equal(expect))
			},
			Entry(
				"default",
				Config{},
				"use ssh agent",
			),
			Entry(
				"explicit key",
				Config{
					SSHKeyFile: "/path/to/key.pem",
				},
				"use ssh key (key.pem)",
			),
			Entry(
				"prefer HTTP",
				Config{
					PreferHTTP: true,
				},
				"use ssh agent, prefer http",
			),
		)
	})
})
