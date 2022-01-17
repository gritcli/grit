package configtest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/internal/daemon/internal/config"
	"github.com/gritcli/grit/internal/daemon/internal/registry"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
)

// SourceDriverTest is a test that tests source driver configuration.
type SourceDriverTest = table.TableEntry

// TestSourceDriver runs a series of tests
func TestSourceDriver(
	r sourcedriver.Registration,
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

			dir, cleanup := makeConfigDir(content)
			defer cleanup()

			cfg, err := config.Load(dir, reg)
			expect(dir, cfg, err)
		},
		tests...,
	)
}

// Success returns a test that tests a source driver configuration that is
// expected to pass.
func Success(
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

// Failure returns a test that tests a source driver configuration that is
// expected to fail.
func Failure(
	description string,
	content string,
	expect string,
) SourceDriverTest {
	return table.Entry(
		description,
		content,
		func(dir string, cfg config.Config, err error) {
			message := strings.ReplaceAll(err.Error(), dir, "<dir>")
			gomega.Expect(message).To(gomega.Equal(expect))
		},
	)
}

// makeConfigDir makes a temporary config directory containing config files
// containing the given configuration content.
func makeConfigDir(configs ...string) (dir string, cleanup func()) {
	dir, err := os.MkdirTemp("", "")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	for i, cfg := range configs {
		err := os.WriteFile(
			filepath.Join(dir, fmt.Sprintf("config-%d.hcl", i)),
			[]byte(cfg),
			0600,
		)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	}

	return dir, func() {
		os.RemoveAll(dir)
	}
}
