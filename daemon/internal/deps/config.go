package deps

import (
	"os"

	"github.com/gritcli/grit/builtins/githubsource"
	"github.com/gritcli/grit/builtins/gitvcs"
	"github.com/gritcli/grit/config"
)

func init() {
	Container.Provide(func() *config.DriverRegistry {
		r := &config.DriverRegistry{}
		r.RegisterSourceDriver("github", githubsource.Registration)
		r.RegisterVCSDriver("git", gitvcs.Registration)

		return r
	})

	Container.Provide(func(r *config.DriverRegistry) (config.Config, error) {
		return config.Load(
			os.Getenv("GRIT_CONFIG_DIR"),
			r,
		)
	})
}
