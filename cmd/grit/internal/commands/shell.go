package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

// newSHellIntegrationCommand returns the "shell-integration" command.
func newShellIntegrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell-integration",
		Short: "Setup shell integration",
		RunE: func(
			cmd *cobra.Command,
			args []string,
		) error {
			return errors.New("not implemented")
		},
	}
}
