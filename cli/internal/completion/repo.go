package completion

import (
	"context"

	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/spf13/cobra"
)

// ValidArgsFunc is a function for use with cobra.Command.ValidArgsFunction.
type ValidArgsFunc func(
	*cobra.Command,
	[]string,
	string,
) ([]string, cobra.ShellCompDirective)

// RepoName returns a ValidArgsFunc that completes the argument at the given
// position using the known repository names.
func RepoName(pos int, loc ...api.Locality) ValidArgsFunc {
	return callSuggestAPI(
		pos,
		func(
			ctx context.Context,
			client api.APIClient,
			word string,
		) (*api.SuggestResponse, error) {
			return client.SuggestRepos(
				ctx,
				&api.SuggestReposRequest{
					Word:           word,
					LocalityFilter: loc,
				},
			)
		},
	)
}

// suggestFunc is a function that invokes one of the suggest API operations.
type suggestFunc func(
	ctx context.Context,
	client api.APIClient,
	word string,
) (*api.SuggestResponse, error)

// callSuggestAPI returns a ValidArgsFunc that completes the argument at the
// given position by calling fn().
func callSuggestAPI(pos int, fn suggestFunc) ValidArgsFunc {
	return func(
		cmd *cobra.Command,
		args []string,
		word string,
	) (words []string, _ cobra.ShellCompDirective) {
		if len(args) != pos {
			return nil, cobra.ShellCompDirectiveDefault
		}

		err := cobradi.Invoke(cmd, func(
			ctx context.Context,
			client api.APIClient,
		) error {
			res, err := fn(ctx, client, word)
			words = res.GetWords()
			return err
		})

		if err != nil {
			cobra.CompErrorln(err.Error())
		}

		return words, cobra.ShellCompDirectiveNoFileComp
	}
}
