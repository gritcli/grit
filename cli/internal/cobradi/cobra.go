package cobradi

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

// Invoke invokes fn with arguments populated by the container associated with
// the given command.
func Invoke(cmd *cobra.Command, fn interface{}) error {
	ctx := cmd.Context()
	con := ctx.Value(contextKey{}).(*di.Container)

	con.Provide(func() context.Context {
		return ctx
	})

	return con.Invoke(fn)
}

// Provide registers fn as a provider function for the container. Its return
// values are added to the container.
func Provide(cmd *cobra.Command, fn interface{}) {
	con := cmd.Context().Value(contextKey{}).(*di.Container)
	con.Provide(fn)
}
