package flags

import (
	"github.com/gritcli/grit/common/config"
	"github.com/spf13/cobra"
)

// SetupConfig sets up --config flag on the root command.
func SetupConfig(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		"config", "c",
		config.DefaultDirectory,
		"set the path to the Grit configuration directory",
	)
}

// ConfigPath returns the path to the configuration directory based on the
// --config flag.
func ConfigPath(cmd *cobra.Command) string {
	dir, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}

	return dir
}
