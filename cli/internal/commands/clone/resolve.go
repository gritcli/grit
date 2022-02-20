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
) (*api.RemoteRepo, bool, error) {
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
		return nil, false, err
	}

	resolver := resolveNonInteractive
	if flags.IsInteractive(cmd) {
		resolver = resolveInteractive
	}

	repo, ok, err := resolver(ctx, cmd, query, responses)
	if err != nil {
		return nil, false, err
	}

	return repo, ok, nil
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
) (*api.RemoteRepo, bool, error) {
	p := tea.NewProgram(
		newResolveModel(
			query,
			responses,
		),
		tea.WithInput(cmd.InOrStdin()),
		tea.WithOutput(cmd.OutOrStdout()),
	)

	x, err := p.StartReturningModel()
	if err != nil {
		return nil, false, err
	}

	m := x.(resolveModel)
	return m.Repo, m.Repo != nil, m.Error
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
) (*api.RemoteRepo, bool, error) {
	var repos []*api.RemoteRepo

	for {
		res, err := responses.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, false, err
		}

		if out := res.GetOutput(); out != nil {
			cmd.Println(out.GetMessage())
		} else if r := res.GetRemoteRepo(); r != nil {
			repos = append(repos, r)
		}
	}

	switch len(repos) {
	case 1:
		return repos[0], true, nil
	case 0:
		return nil, false, errors.New("no matching repositories")
	default:
		printUnambiguousCommands(cmd, query, repos)
		return nil, false, errors.New("multiple matching repositories")
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

	socket, _ := flags.Socket(cmd)

	for n, r := range repos {
		cmd.PrintErrf("%d) %s\n", n+1, r.GetName())
		cmd.Println("")
		cmd.PrintErrf(
			"  %s clone \\\n",
			os.Args[0],
		)

		if socket != api.DefaultSocket {
			cmd.PrintErrf("    --socket %s \\\n", shell.Escape(socket))
		}

		cmd.PrintErrf("    --from-source %s \\\n", shell.Escape(r.GetSource()))
		cmd.PrintErrf("    --no-resolve %s\n", shell.Escape(r.GetId()))
		cmd.Println("")
	}
}
