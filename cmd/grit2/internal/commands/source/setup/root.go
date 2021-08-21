package setup

import (
	"github.com/spf13/cobra"
)

// NewRoot returns the "source setup" command.
func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "add a configure a repository source",
	}

	cmd.AddCommand(
		newGitHubCommand(),
	)

	return cmd
}
