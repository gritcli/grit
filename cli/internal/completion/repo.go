package completion

import (
	"context"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/spf13/cobra"
)

// RepoName returns a ValidArgsFunc that completes the argument at the given
// position using the known repository names.
func RepoName(
	con *imbue.Container,
	loc ...api.Locality,
) ValidArgsFunc {
	return useAPI(
		con,
		func(
			ctx context.Context,
			client api.APIClient,
			cmd *cobra.Command,
			args []string,
			word string,
		) ([]string, cobra.ShellCompDirective, error) {
			res, err := client.SuggestRepos(
				ctx,
				&api.SuggestReposRequest{
					Word:           word,
					LocalityFilter: loc,
				},
			)
			return res.GetWords(), cobra.ShellCompDirectiveNoFileComp, err
		},
	)
}
