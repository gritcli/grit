package configtest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/onsi/gomega"
)

// writeConfigs makes a temporary config directory containing config files
// containing the given configuration content.
func writeConfigs(configs ...string) (dir string, cleanup func()) {
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
