package deps

import (
	"os"

	"github.com/gritcli/grit/internal/config"
)

func init() {
	Catalog.Provide(func() (config.Config, error) {
		dir := os.Getenv("GRIT_CONFIG_DIR")
		if dir == "" {
			dir = config.DefaultDirectory
		}

		return config.Load(dir)
	})
}
