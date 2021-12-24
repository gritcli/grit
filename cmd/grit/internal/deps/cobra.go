package deps

import (
	"context"

	"github.com/gritcli/grit/internal/di"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

type contextKey struct{}

// Execute executes a Cobra CLI command within the context of c.
//
// It closes the container after the command is completed.
func Execute(ctx context.Context, c *di.Container, cmd *cobra.Command) (err error) {
	ctx = context.WithValue(ctx, contextKey{}, c)

	defer func() {
		err = multierr.Append(
			err,
			c.Close(),
		)
	}()

	return cmd.ExecuteContext(ctx)
}

// Run returns a function that invokes fn with arguments populated by the
// container. The returned function matches the signature of cobra.Command.RunE.
func Run(fn interface{}) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		con := ctx.Value(contextKey{}).(*di.Container)

		return con.Invoke(fn)
	}
}
