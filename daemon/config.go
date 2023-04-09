package daemon

import (
	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/daemon/internal/builtins/githubsource"
	"github.com/gritcli/grit/daemon/internal/builtins/gitvcs"
)

var configDir = ferrite.
	String("GRIT_CONFIG_DIR", "the directory containing Grit's configuration files").
	WithDefault(config.DefaultDirectory).
	Required()

func init() {
	imbue.With0(
		catalog,
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
		catalog,
		func(
			ctx imbue.Context,
			r *config.DriverRegistry,
		) (config.Config, error) {
			return config.Load(
				configDir.Value(),
				r,
			)
		},
	)
}
