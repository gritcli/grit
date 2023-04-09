package signout

import (
	"context"
	_ "embed"
	"errors"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/completion"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/spf13/cobra"
)

//go:embed help.txt
var helpText string

// Command returns the "source auth" command.
func Command(con *imbue.Container) *cobra.Command {
	return &cobra.Command{
		Use:                   "sign-out",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Aliases:               []string{"signout", "logout"},
		Short:                 "Sign out of a repository source",
		Long:                  helpText,
		ValidArgsFunction: completion.Positional(
			completion.SourceName(con),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !flags.IsInteractive(cmd) {
				return errors.New("non-interactive mode is not supported")
			}

			cmd.SilenceUsage = true

			return imbue.Invoke2(
				cmd.Context(),
				con,
				func(
					ctx context.Context,
					client api.APIClient,
					clientOptions *api.ClientOptions,
				) error {
					ctx, cancel := context.WithCancel(ctx)
					defer cancel()

					req := &api.SignOutRequest{
						Source: args[0],
					}

					_, err := client.SignOut(ctx, req)
					if err != nil {
						return err
					}

					return nil
				},
			)
		},
	}
}
