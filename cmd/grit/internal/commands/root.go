package commands

import (
	"context"
	"io"
	"os"

	"github.com/gritcli/grit/cmd/grit/internal/commands/source"
	"github.com/gritcli/grit/cmd/grit/internal/deps"
	"github.com/gritcli/grit/cmd/grit/internal/flags"
	"github.com/gritcli/grit/cmd/grit/internal/shell"
	"github.com/gritcli/grit/internal/config"
	"github.com/gritcli/grit/internal/di"
	"github.com/spf13/cobra"
)

// NewRoot returns the root command.
//
// v is the version to display. It is passed from the main package where it is
// made available as part of the build process.
func NewRoot(v string) *cobra.Command {
	root := &cobra.Command{
		Version:      v,
		Use:          "grit",
		Short:        "Manage your local source repository clones",
		SilenceUsage: true, // otherwise ANY error shows usage
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Add the currently-executing Cobra CLI command the the DI
			// container. This hook is called after the CLI arguments are
			// resolved to a specific command.
			deps.Container.Provide(func() (
				context.Context,
				*cobra.Command,
				[]string,
			) {
				return cmd.Context(), cmd, args
			})
		},
	}

	flags.SetupVerbose(root)

	provideConfig(&deps.Container, root)
	provideShellExecutor(&deps.Container, root)

	root.AddCommand(
		source.NewRoot(),
		newShellIntegrationCommand(),
		newCloneCommand(),
	)

	return root
}

// provideConfig sets up the DI container to supply config.Config values based
// on the path specified by the --config flag.
func provideConfig(c *di.Container, root *cobra.Command) {
	// Add --config as a "persistent" flag so that it's available on all
	// commands.
	root.PersistentFlags().StringP(
		"config", "c",
		config.DefaultDirectory,
		"set the path to the Grit configuration directory",
	)

	// Setup a DI provider that provides config.Config value using the --config
	// flag to determine where to load the config from.
	c.Provide(func(cmd *cobra.Command) (config.Config, error) {
		dir, err := cmd.Flags().GetString("config")
		if err != nil {
			return config.Config{}, err
		}

		cfg, err := config.Load(dir)
		if err != nil {
			return config.Config{}, err
		}

		return cfg, nil
	})
}

// provideShellExecutor sets up the DI container to supply a shell.Executor that
// writes to the file specified by the --shell-executor-output flag.
func provideShellExecutor(c *di.Container, root *cobra.Command) {
	// Add --shell-executor-output as a "persistent" flag so it's available on
	// all commands. Mark it as hidden as it should only be passed by the code
	// generated by the "shell-integration" command, and never by the user
	// directly.
	f := root.PersistentFlags()
	f.String("shell-executor-output", "", "output file for shell commands to execute")
	f.MarkHidden("shell-executor-output") //nolint:errcheck

	// Setup a DI provider that provides a shell.Executor that writes to the
	// file specified by the --shell-executor-output flag.
	c.Provide(func(
		cmd *cobra.Command,
	) (shell.Executor, error) {
		filename, err := cmd.Flags().GetString("shell-executor-output")
		if err != nil {
			return nil, err
		}

		if filename == "" {
			cmd.PrintErrf("Shell integration has not been configured. For more information run:\n\n")
			cmd.PrintErrf("    %s help shell-integration\n\n", os.Args[0])
			return shell.NewExecutor(io.Discard), nil
		}

		fp, err := os.Create(filename)
		if err != nil {
			return nil, err
		}

		c.Defer(func() error {
			defer fp.Close()

			if err := fp.Sync(); err != nil {
				return err
			}

			return fp.Close()
		})

		return shell.NewExecutor(fp), nil
	})

}
