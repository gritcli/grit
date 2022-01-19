package deps

import (
	"os"

	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver/githubsource"
	"github.com/gritcli/grit/driver/vcsdriver/gitvcs"
)

func init() {
	Container.Provide(func() *registry.Registry {
		r := &registry.Registry{}
		r.RegisterSourceDriver("github", githubsource.Registration)
		r.RegisterVCSDriver("git", gitvcs.Registration)

		return r
	})

	Container.Provide(func(r *registry.Registry) (config.Config, error) {
		return config.Load(
			os.Getenv("GRIT_CONFIG_DIR"),
			r,
		)
	})
}
