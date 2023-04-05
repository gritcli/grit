package completion

import (
	"context"
	"strings"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/spf13/cobra"
)

// SourceName returns a ValidArgsFunc that completes the argument at the given
// position using the known repository source names.
func SourceName(
	con *imbue.Container,
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
			res, err := client.ListSources(
				ctx,
				&api.ListSourcesRequest{},
			)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp, err
			}

			var words []string
			for _, src := range res.GetSources() {
				name := src.GetName()

				if strings.HasPrefix(name, word) {
					words = append(words, name)
				}
			}

			return words, cobra.ShellCompDirectiveNoFileComp, nil
		},
	)
}
