package deps

import (
	"os"

	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/registry"
)

func init() {
	Container.Provide(func(r *registry.Registry) (config.Config, error) {
		return config.Load(
			os.Getenv("GRIT_CONFIG_DIR"),
			&registry.Registry{
				Parent: &registry.BuiltIns,
			},
		)
	})
}
