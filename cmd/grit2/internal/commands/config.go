package commands

import (
	"os"

	"github.com/gritcli/grit/internal/config"
	"github.com/gritcli/grit/internal/di"
	"github.com/spf13/cobra"
)

// setupConfig adds the --config flag as a persistent flag on the root
// command so that it is available to all commands.
func setupConfig(root *cobra.Command) {
	root.PersistentFlags().StringP(
		"config", "c",
		config.DefaultFile,
		"set the path to the Grit configuration file",
	)
}

// provideConfig adds parses the Grit configuration and adds the config.Config
// to the DI configuration.
func provideConfig(c *di.Container, cmd *cobra.Command) {
	c.Provide(func() (config.Config, error) {
		filename, err := cmd.Flags().GetString("config")
		if err != nil {
			return config.Config{}, err
		}

		cfg, err := config.ParseFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				if !cmd.Flags().Changed("config") {
					// If the --config flag was not specified and the config
					// file doesn't exist we fall back to the default config.
					return config.DefaultConfig, nil
				}
			}

			return config.Config{}, err
		}

		return cfg, nil
	})
}
