package configtest

import (
	"strings"

	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// SourceDriverTest is a test that tests source driver configuration.
type SourceDriverTest = table.TableEntry

// TestSourceDriver runs a series of tests
func TestSourceDriver(
	r sourcedriver.Registration,
	zero sourcedriver.Config,
	deps []vcsdriver.Registration,
	tests ...SourceDriverTest,
) {
	table.DescribeTable(
		"it loads the configuration",
		func(
			content string,
			expect func(dir string, cfg config.Config, err error),
		) {
			reg := &registry.Registry{}
			reg.RegisterSourceDriver(r.Name, r)

			for _, d := range deps {
				reg.RegisterVCSDriver(d.Name, d)
			}

			dir, cleanup := writeConfigs(content)
			defer cleanup()

			cfg, err := config.Load(dir, reg)
			expect(dir, cfg, err)
		},
		tests...,
	)
}

// SourceSuccess returns a test that tests a source driver configuration that is
// expected to pass.
func SourceSuccess(
	description string,
	content string,
	expect sourcedriver.Config,
) SourceDriverTest {
	return table.Entry(
		description,
		content,
		func(dir string, cfg config.Config, err error) {
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(cfg.Sources).To(gomega.HaveLen(1))
			gomega.Expect(cfg.Sources[0].Driver).To(gomega.Equal(expect))
		},
	)
}

// SourceFailure returns a test that tests a source driver configuration that is
// expected to fail.
func SourceFailure(
	description string,
	content string,
	expect string,
) SourceDriverTest {
	return table.Entry(
		description,
		content,
		func(dir string, cfg config.Config, err error) {
			orig := format.TruncatedDiff
			format.TruncatedDiff = false
			defer func() {
				format.TruncatedDiff = orig
			}()

			message := strings.ReplaceAll(err.Error(), dir, "<dir>")
			gomega.Expect(message).To(gomega.Equal(expect))
		},
	)
}
