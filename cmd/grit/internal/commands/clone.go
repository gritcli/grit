package commands

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/cmd/grit/internal/deps"
	"github.com/gritcli/grit/internal/api"
	"github.com/spf13/cobra"
)

// newCloneCommand returns the "clone" command.
func newCloneCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <repo>",
		Args:  cobra.ExactArgs(1),
		Short: "Clone a remote repository",
		Long: heredoc.Doc(`
		The "clone" command makes a local clone of a remote repository then
		changes the shell's current working directory to the clone's working
		tree.

		The <repo> argument is a repository name (or a part thereof), URL, or
		other identifier. For example, the Grit repository itself may be
		referred to as "gritcli/grit", just "grit".

		Each of the repository sources defined in the Grit configuration file is
		searched for matches to the provided repository name. If there are
		multiple matches and the shell is interactive the user is prompted to
		select the desired repository.
		`),
		RunE: deps.Run(func(
			ctx context.Context,
			cmd *cobra.Command,
			args []string,
			client api.APIClient,
			clientOptions *api.ClientOptions,
		) error {
			if args[0] == "" {
				return errors.New("<repo> argument must not be empty")
			}

			req := &api.ResolveRequest{
				ClientOptions: clientOptions,
				Query:         args[0],
			}

			stream, err := client.Resolve(ctx, req)
			if err != nil {
				return err
			}

			for {
				res, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}

					return err
				}

				fmt.Println(res)
			}

			return nil
		}),
	}
}
