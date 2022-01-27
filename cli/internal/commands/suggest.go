package commands

import (
	"context"

	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/deps"
	"github.com/spf13/cobra"
)

// suggestFunc is a function that returns a suggestion response from the API.
type suggestFunc func(
	ctx context.Context,
	client api.APIClient,
	cmd *cobra.Command,
	args []string,
	word string,
) (*api.SuggestResponse, error)

// validArgsFunc is a function for use with cobra.Command.ValidArgsFunction.
type validArgsFunc func(
	*cobra.Command,
	[]string,
	string,
) ([]string, cobra.ShellCompDirective)

// suggest is a helper function for using the suggestion API for autocompletion.
func suggest(fn suggestFunc) validArgsFunc {
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
