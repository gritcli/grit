package add

import (
	"fmt"
	"os"

	"github.com/cli/oauth"
	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/spf13/cobra"
)

// newGitHubCommand returns the "source add github" command.
func newGitHubCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "github",
		Short: "add github or github enterprise as a repository source",
		RunE: di.RunE(func(
			cmd *cobra.Command,
			args []string,
		) error {
			flow := &oauth.Flow{
				Hostname:     "github.com",
				ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
				ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"), // only applicable to web app flow
				CallbackURI:  "http://127.0.0.1/callback",      // only applicable to web app flow
				Scopes:       []string{"repo", "read:org", "gist"},
			}

			accessToken, err := flow.DetectFlow()
			if err != nil {
				panic(err)
			}

			fmt.Printf("Access token: %s\n", accessToken.Token)
		}),
	}

	cmd.Flags().String(
		"enterprise",
		"",
		"set the hostname of a GitHub Enterprise server hostname",
	)

	return cmd
}
