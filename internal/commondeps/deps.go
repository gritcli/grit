package commondeps

import (
	"os"
	"path/filepath"

	"github.com/gritcli/grit/internal/di"
)

// ExecutableInfo contains information about the current executable.
type ExecutableInfo struct {
	Name    string
	Version string
}

// Provide adds a provider to c that provides an ExecutableInfo
// value.
func Provide(c *di.Container, version string) {
	c.Provide(func() ExecutableInfo {
		return ExecutableInfo{
			filepath.Base(os.Args[0]),
			version,
		}
	})
}
