package setup

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cli/oauth"
	"github.com/gritcli/grit/internal/di/cobradi"
	"github.com/gritcli/grit/internal/source/githubsource"
	"github.com/spf13/cobra"
)

// newGitHubCommand returns the "source setup github" command.
func newGitHubCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "github",
		Short: "setup GitHub as a repository source",
		Long: heredoc.Doc(`
		The "github" command configures Grit to clone repositories from a GitHub
		installation.

		By default it configures github.com as a repository source. The --enterprise
		flag can be used to specify the address (<hostname>:<port>) of a GitHub
		Enterprise Cloud or GitHub Enterprise Server API server instead.
		`),
		RunE: cobradi.RunE(func(
			cmd *cobra.Command,
		) error {
			if addr, _ := cmd.Flags().GetString("enterprise"); addr != "" {
				return setupGitHubEnterprise(cmd, addr)
			}

			return setupGitHubDotCom(cmd)
		}),
	}

	cmd.Flags().String(
		"enterprise",
		"",
		"set the hostname (and optionally port) of a GitHub Enterprise installation",
	)

	cmd.Flags().String(
		"auth-token",
		"",
		"use a GitHub personal access token for authentication",
	)

	return cmd
}

// setupGitHubDotCom configures github.com as a repository source.
func setupGitHubDotCom(cmd *cobra.Command) error {
	flow := &oauth.Flow{
		Hostname: "github.com",
		ClientID: githubsource.AppClientID,
		Scopes:   githubsource.RequiredScopes,
	}

	token, err := flow.DeviceFlow()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Access token: %s\n", token.Token)

	return nil
}

// setupGitHubEnterprise configures a GitHub Enterprise server as a repository
// source.
func setupGitHubEnterprise(cmd *cobra.Command, addr string) error {
	panic("not implemented")
}
