package configtest

import (
	"strings"

	"github.com/gritcli/grit/daemon/internal/config"
	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// SourceDriverTest is a test that tests source driver configuration.
type SourceDriverTest = ginkgo.TableEntry

// TestSourceDriver runs a series of tests
func TestSourceDriver(
	r sourcedriver.Registration,
	zero sourcedriver.Config,
	deps []vcsdriver.Registration,
	tests ...SourceDriverTest,
) {
	args := []any{
		func(
			content string,
			expect func(dir string, cfg config.Config, err error),
		) {
			reg := &config.DriverRegistry{}
			reg.RegisterSourceDriver(r.Name, r)

			for _, d := range deps {
				reg.RegisterVCSDriver(d.Name, d)
			}

			dir, cleanup := writeConfigs(content)
			defer cleanup()

			cfg, err := config.Load(dir, reg)
			expect(dir, cfg, err)
		},
	}
	for _, test := range tests {
		args = append(args, test)
	}

	ginkgo.DescribeTable("it loads the configuration", args...)
}

// SourceSuccess returns a test that tests a source driver configuration that is
// expected to pass.
func SourceSuccess(
	description string,
	content string,
	expect sourcedriver.Config,
) SourceDriverTest {
	return ginkgo.Entry(
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
	return ginkgo.Entry(
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
