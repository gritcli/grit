package deps

import (
	"os"

	"github.com/gritcli/grit/common/config"
)

func init() {
	Container.Provide(func() (config.Config, error) {
		dir := os.Getenv("GRIT_CONFIG_DIR")
		if dir == "" {
			dir = config.DefaultDirectory
		}

		return config.Load(dir)
	})
}