package clone

import (
	"context"
	"errors"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

// resolve resolves a repository query to a single repository.
func resolve(
	ctx context.Context,
	cmd *cobra.Command,
	client api.APIClient,
	clientOptions *api.ClientOptions,
	query string,
	source string,
) (string, string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req := &api.ResolveRepoRequest{
		ClientOptions: clientOptions,
		Query:         query,
		LocalityFilter: []api.Locality{
			api.Locality_REMOTE,
		},
	}

	if source != "" {
		req.SourceFilter = []string{source}
	}

	responses, err := client.ResolveRepo(ctx, req)
	if err != nil {
		return "", "", err
	}

	resolver := resolveNonInteractive
	if flags.IsInteractive(cmd) {
		resolver = resolveInteractive
	}

	repo, err := resolver(ctx, cmd, query, responses)
	if err != nil {
		return "", "", err
	}

	return repo.GetId(), repo.GetSource(), nil
}

// resolveInteractive resolves the repository query to a single repository.
//
// If the query is ambiguous the user is prompted to choose from the matching
// repositories.
func resolveInteractive(
	ctx context.Context,
	cmd *cobra.Command,
	query string,
	responses api.API_ResolveRepoClient,
) (*api.RemoteRepo, error) {
	p := tea.NewProgram(newResolveModel(
		query,
		responses,
	))

	x, err := p.StartReturningModel()
	if err != nil {
		return nil, err
	}

	m := x.(resolveModel)
	return m.Repo, m.Error
}

// resolveNonInteractive resolves the repository query to a single repository ID
// without allowing user interaction.
//
// If the query is ambiguous it prints instructions for invoking the command
// with an unambiguous query then returns an error.
func resolveNonInteractive(
	ctx context.Context,
	cmd *cobra.Command,
	query string,
	responses api.API_ResolveRepoClient,
) (*api.RemoteRepo, error) {
	var repos []*api.RemoteRepo

	for {
		res, err := responses.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if out := res.GetOutput(); out != nil {
			cmd.Println(out.GetMessage())
		} else if r := res.GetRemoteRepo(); r != nil {
			repos = append(repos, r)
		}
	}

	switch len(repos) {
	case 1:
		return repos[0], nil
	case 0:
		return nil, errors.New("no matching repositories")
	default:
		printUnambiguousCommands(cmd, query, repos)
		return nil, errors.New("multiple matching repositories")
	}
}

// printUnambigousCommands prints the unambigous commands that can be run when
// the repository query matches multiple repositories.
func printUnambiguousCommands(
	cmd *cobra.Command,
	query string,
	repos []*api.RemoteRepo,
) {
	cmd.PrintErrf(
		"Multiple repositories match '%s'. Use one of the following commands instead:\n\n",
		query,
	)

	for _, r := range repos {
		cmd.PrintErrf(
			"  %s clone --from-source %s --no-resolve %s # %s \n",
			os.Args[0],
			shell.Escape(r.GetSource()),
			shell.Escape(r.GetId()),
			r.GetName(),
		)
	}

	cmd.Println("")
}
