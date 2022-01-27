package commands

import (
	"context"
	"errors"
	"io"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/interactive"
	"github.com/gritcli/grit/cli/internal/render"
	"github.com/gritcli/grit/cli/internal/shell"
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
		changes the users's current working directory to that of the clone.

		The <repo> argument is a repository name (or a part thereof), URL, or
		other identifier. For example, the Grit repository itself may be
		referred to as "gritcli/grit" or simply "grit".

		Each of the repository sources defined in the Grit configuration file is
		searched for matches to the provided repository name. If there are
		multiple matches and the shell is interactive the user is prompted to
		select the desired repository.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("<repo> argument must not be empty")
			}

			return cobradi.Invoke(cmd, func(
				ctx context.Context,
				client api.APIClient,
				clientOptions *api.ClientOptions,
				executor shell.Executor,
			) error {
				repo, err := resolveRepo(
					ctx,
					cmd,
					client,
					clientOptions,
					args[0],
				)
				if err != nil {
					return err
				}

				dir, err := cloneRepo(
					ctx,
					cmd,
					client,
					clientOptions,
					repo,
				)
				if err != nil {
					return err
				}

				cmd.Println(render.RelPath(dir))

				return executor("cd", dir)
			})
		},
		ValidArgsFunction: suggest(func(
			ctx context.Context,
			client api.APIClient,
			cmd *cobra.Command,
			args []string,
			word string,
		) (*api.SuggestResponse, error) {
			if len(args) != 0 {
				return nil, nil
			}

			return client.SuggestRepo(ctx, &api.SuggestRepoRequest{
				Word:          word,
				IncludeLocal:  false,
				IncludeRemote: true,
			})
		}),
	}
}

// resolveRepo resolves a repo query string to a specific repo.
func resolveRepo(
	ctx context.Context,
	cmd *cobra.Command,
	client api.APIClient,
	clientOptions *api.ClientOptions,
	query string,
) (*api.Repo, error) {
	req := &api.ResolveRequest{
		ClientOptions: clientOptions,
		Query:         query,
	}

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

	return interactive.SelectRepos(
		cmd,
		"Which repository would you like to clone?",
		repos,
	)
}

// cloneRepo clones a repository.
func cloneRepo(
	ctx context.Context,
	cmd *cobra.Command,
	client api.APIClient,
	clientOptions *api.ClientOptions,
	repo *api.Repo,
) (string, error) {
	req := &api.CloneRequest{
		ClientOptions: clientOptions,
		Source:        repo.Source,
		RepoId:        repo.Id,
	}

	stream, err := client.Clone(ctx, req)
	if err != nil {
		return "", err
	}

	dir := ""

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		if out := res.GetOutput(); out != nil {
			cmd.Println(out.Message)
		} else if d := res.GetDirectory(); d != "" {
			dir = d
		}
	}

	if dir == "" {
		return "", errors.New("server did not provide the directory of the clone")
	}

	return dir, nil
}
