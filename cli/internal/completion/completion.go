package completion

import (
	"context"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/api"
	"github.com/spf13/cobra"
)

// ValidArgsFunc is a function for use with cobra.Command.ValidArgsFunction.
type ValidArgsFunc func(
	*cobra.Command,
	[]string,
	string,
) ([]string, cobra.ShellCompDirective)

// Positional returns a ValidArgsFunc that delegates to the n'th function when
// the n'th argument is being completed.
func Positional(funcs ...ValidArgsFunc) ValidArgsFunc {
	return func(
		cmd *cobra.Command,
		args []string,
		word string,
	) (words []string, _ cobra.ShellCompDirective) {
		if len(args) >= len(funcs) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return funcs[len(args)](cmd, args, word)
	}
}

func useAPI(
	con *imbue.Container,
	fn func(
		ctx context.Context,
		client api.APIClient,
		cmd *cobra.Command,
		args []string,
		word string,
	) ([]string, cobra.ShellCompDirective, error),
) ValidArgsFunc {
	return func(
		cmd *cobra.Command,
		args []string,
		word string,
	) (words []string, dir cobra.ShellCompDirective) {
		err := imbue.Invoke1(
			cmd.Context(),
			con,
			func(
				ctx context.Context,
				client api.APIClient,
			) error {
				var err error
				words, dir, err = fn(ctx, client, cmd, args, word)
				return err
			},
		)

		if err != nil {
			cobra.CompErrorln(err.Error())
		}

		return words, dir
	}
}
