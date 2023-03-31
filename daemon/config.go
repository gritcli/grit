package daemon

import (
	"os"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/builtins/githubsource"
	"github.com/gritcli/grit/builtins/gitvcs"
	"github.com/gritcli/grit/config"
)

func init() {
	imbue.With0(
		container,
		func(
			ctx imbue.Context,
		) (*config.DriverRegistry, error) {
			r := &config.DriverRegistry{}
			r.RegisterSourceDriver("github", githubsource.Registration)
			r.RegisterVCSDriver("git", gitvcs.Registration)
			return r, nil
		},
	)

	imbue.With1(
		container,
		func(
			ctx imbue.Context,
			r *config.DriverRegistry,
		) (config.Config, error) {
			return config.Load(
				os.Getenv("GRIT_CONFIG_DIR"),
				r,
			)
		},
	)
}
