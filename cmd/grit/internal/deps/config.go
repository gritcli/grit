package deps

import (
	"github.com/gritcli/grit/cmd/grit/internal/flags"
	"github.com/gritcli/grit/common/config"
	"github.com/spf13/cobra"
)

func init() {
	Container.Provide(func(
		cmd *cobra.Command,
	) (config.Config, error) {
		return config.Load(
			flags.ConfigPath(cmd),
		)
	})
}
