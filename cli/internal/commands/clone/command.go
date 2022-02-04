package clone

import (
	"context"
	"errors"
	"io"

	"github.com/MakeNowJust/heredoc/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/render"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

// NewCommand returns the "clone" command.
func NewCommand() *cobra.Command {
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

		If there are multiple matching local clones and the shell is interactive
		the user is prompted to select the desired repository.
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
				streamCtx, streamCancel := context.WithCancel(ctx)
				defer streamCancel()

				stream, err := client.ResolveRepo(streamCtx, &api.ResolveRepoRequest{
					ClientOptions: clientOptions,
					Query:         args[0],
					Locality:      api.Locality_REMOTE_ONLY,
				})
				if err != nil {
					return err
				}

				p := tea.NewProgram(model{
					Stream: stream,
				})

				x, err := p.StartReturningModel()
				if err != nil {
					return err
				}

				streamCancel()

				m := x.(model)

				if m.ChosenRepo == nil {
					return m.Error
				}

				local, err := cloneRepo(
					ctx,
					cmd,
					client,
					clientOptions,
					m.ChosenRepo,
				)
				if err != nil {
					return err
				}

				cmd.Println(render.RelPath(local.AbsoluteCloneDir))

				return executor("cd", local.AbsoluteCloneDir)
			})
		},
		// ValidArgsFunction: suggest(func(
		// 	ctx context.Context,
		// 	client api.APIClient,
		// 	cmd *cobra.Command,
		// 	args []string,
		// 	word string,
		// ) (*api.SuggestResponse, error) {
		// 	if len(args) != 0 {
		// 		return nil, nil
		// 	}

		// 	return client.SuggestRepos(ctx, &api.SuggestReposRequest{
		// 		Word:     word,
		// 		Locality: api.Locality_REMOTE_ONLY,
		// 	})
		// }),
	}
}

// cloneRepo clones a remote repository using the clone API operation.
func cloneRepo(
	ctx context.Context,
	cmd *cobra.Command,
	client api.APIClient,
	clientOptions *api.ClientOptions,
	repo *api.RemoteRepo,
) (*api.LocalRepo, error) {
	req := &api.CloneRepoRequest{
		ClientOptions: clientOptions,
		Source:        repo.Source,
		RepoId:        repo.Id,
	}

	stream, err := client.CloneRepo(ctx, req)
	if err != nil {
		return nil, err
	}

	var local *api.LocalRepo

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
		} else if r := res.GetLocalRepo(); r != nil {
			local = r
		}
	}

	if local == nil {
		return nil, errors.New("server did not provide information about the local clone")
	}

	return local, nil
}
