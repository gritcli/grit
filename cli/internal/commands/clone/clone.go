package clone

import (
	"context"
	"errors"
	"io"

	"github.com/gritcli/grit/api"
	"github.com/spf13/cobra"
)

// clone makes a local clone of a remote repository.
func clone(
	ctx context.Context,
	cmd *cobra.Command,
	client api.APIClient,
	options *api.ClientOptions,
	id, source string,
) (string, error) {
	req := &api.CloneRepoRequest{
		ClientOptions: options,
		Source:        source,
		RepoId:        id,
	}

	responses, err := client.CloneRepo(ctx, req)
	if err != nil {
		return "", err
	}

	var local *api.LocalRepo

	for {
		res, err := responses.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		if out := res.GetOutput(); out != nil {
			cmd.Println(out.Message)
		} else if r := res.GetLocalRepo(); r != nil {
			local = r
		}
	}

	if local == nil {
		return "", errors.New("server did not provide information about the local clone")
	}

	return local.GetAbsoluteCloneDir(), nil
}
