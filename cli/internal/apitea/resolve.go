package apitea

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gritcli/grit/api"
)

// RemoteRepoMatched is a tea.Msg that indicates a ResolveRepo() API operation
// has found a matching remote repository.
type RemoteRepoMatched struct {
	Repo *api.RemoteRepo
	Tag  string
}

// LocalRepoMatched is a tea.Msg that indicates a ResolveRepo() API operation
// has found a matching local repository.
type LocalRepoMatched struct {
	Repo *api.LocalRepo
	Tag  string
}

// ResolveComplete is a tea.Msg that indicates a ResolveRepo() API operation has
// completed successfully.
type ResolveComplete struct {
	Tag string
}

// ResolveFailed is a tea.Msg that indicates a ResolveRepo() API operation has
// failed.
type ResolveFailed struct {
	Error error
	Tag   string
}

// WaitForResolveResponse waits for the next response of a ResolveRepo() API
// operation and dispatches the appropriate message.
func WaitForResolveResponse(
	responses api.API_ResolveRepoClient,
	tag string,
) tea.Cmd {
	return func() tea.Msg {
		res, err := responses.Recv()

		if err == io.EOF {
			return ResolveComplete{
				Tag: tag,
			}
		}

		if err != nil {
			return ResolveFailed{
				Error: err,
				Tag:   tag,
			}
		}

		if out := res.GetOutput(); out != nil {
			return OutputReceived{
				Output: out,
				Tag:    tag,
			}
		}

		if repo := res.GetRemoteRepo(); repo != nil {
			return RemoteRepoMatched{
				Repo: repo,
				Tag:  tag,
			}
		}

		if repo := res.GetLocalRepo(); repo != nil {
			return LocalRepoMatched{
				Repo: repo,
				Tag:  tag,
			}
		}

		return nil
	}
}
