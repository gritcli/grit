package deps

import (
	"os"

	"github.com/gritcli/grit/internal/daemon/internal/config"
	"github.com/gritcli/grit/internal/daemon/internal/registry"
)

func init() {
	Container.Provide(func(r *registry.Registry) (config.Config, error) {
		dir := os.Getenv("GRIT_CONFIG_DIR")
		if dir == "" {
			dir = config.DefaultDirectory
		}

		return config.Load(
			dir,
			&registry.Registry{
				Parent: &registry.BuiltIns,
			},
		)
	})
}
