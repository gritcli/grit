package add

import (
	"github.com/spf13/cobra"
)

// NewRoot returns the "source add" command.
func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add a new repository source",
	}

	cmd.AddCommand(
		newGitHubCommand(),
	)

	return cmd
}
