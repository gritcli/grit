package commands

import (
	"context"

	"github.com/gritcli/grit/internal/cli/internal/deps"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/spf13/cobra"
)

// suggest is a helper function for using the suggestion API for autocompletion.
//
// It returns a function with a signature that matches
// cobra.Command.ValidArgsFunction.
func suggest(
	fn func(
		ctx context.Context,
		client api.APIClient,
		cmd *cobra.Command,
		args []string,
		word string,
	) (*api.SuggestResponse, error),
) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(
		cmd *cobra.Command,
		args []string,
		word string,
	) ([]string, cobra.ShellCompDirective) {
		var res *api.SuggestResponse

		err := deps.Invoke(cmd, func(
			ctx context.Context,
			client api.APIClient,
		) error {
			var err error
			res, err = fn(ctx, client, cmd, args, word)
			return err
		})

		if err != nil {
			cobra.CompErrorln(err.Error())
		}

		return res.GetWords(), cobra.ShellCompDirectiveNoFileComp
	}
}
