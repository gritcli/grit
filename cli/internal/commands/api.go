package commands

import (
	"context"
	"errors"
	"io"

	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/interactive"
	"github.com/spf13/cobra"
)

// resolveRepo is a helper function for choosing a single repository from the
// result of a call to the resolve API.
func resolveRepo(
	ctx context.Context,
	cmd *cobra.Command,
	client api.APIClient,
	req *api.ResolveRequest,
) (*api.Repo, error) {
	stream, err := client.Resolve(ctx, req)
	if err != nil {
		return nil, err
	}

	var repos []*api.Repo

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if out := res.GetOutput(); out != nil {
			cmd.Println(out.Message)
		} else if r := res.GetRepo(); r != nil {
			repos = append(repos, r)
		}
	}

	if len(repos) == 0 {
		return nil, errors.New("no matching repositories found")
	}

	return interactive.SelectRepos(cmd, repos)
}

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

		err := cobradi.Invoke(cmd, func(
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
