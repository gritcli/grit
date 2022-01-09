package commands

import (
	"errors"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// newSHellIntegrationCommand returns the "shell-integration" command.
func newShellIntegrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell-integration",
		Short: "Setup shell integration",
		Long: heredoc.Doc(`
		This "shell-integration" configures the user's shell to allow Grit to
		change the current working directory.

		Add the following command to your to your shell's initialization script
		(for example .bash_profile, .zprofile, etc).

			eval $(grit shell-integration)

		Note that this does NOT configure shell auto-completion. For
		instructions on configuring auto-completion run:

			grit help completion
		`),
		RunE: func(
			cmd *cobra.Command,
			args []string,
		) error {
			return errors.New("not implemented")
		},
	}
}
