package cli

import (
	"io"
	"os"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/flags"
	"github.com/gritcli/grit/cli/internal/shell"
	"github.com/spf13/cobra"
)

func init() {
	imbue.With1(
		container,
		func(
			ctx imbue.Context,
			cmd *cobra.Command,
		) (shell.Executor, error) {
			filename, ok := flags.ShellExecutorOutputFile(cmd)
			if !ok {
				// Print a warning message if shell integration has not been
				// configured. Without this, grit can not properly change the
				// current working directory of the shell.
				//
				// It's a bit of a hack to do this inside the DI provider, but it
				// means that the warning is only displayed if the command being
				// executed actually requested a shell.Executor.
				if _, ok := flags.ShellExecutorOutputFile(cmd); !ok {
					cmd.PrintErrf("Shell integration has not been configured. For more information run:\n\n")
					cmd.PrintErrf("  %s help setup-shell\n\n", os.Args[0])
				}

				// Return an executor that does nothing, so we can still operate
				// with reduced functionality.
				return shell.NewExecutor(io.Discard), nil
			}

			fp, err := os.Create(filename)
			if err != nil {
				return nil, err
			}

			ctx.Defer(func() error {
				defer fp.Close()

				if err := fp.Sync(); err != nil {
					return err
				}

				return fp.Close()
			})

			return shell.NewExecutor(fp), nil
		},
	)
}
