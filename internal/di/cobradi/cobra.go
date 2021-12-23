package cobradi

import (
	"context"

	"github.com/gritcli/grit/internal/di"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

// Execute executes a command within the context of c.
//
// It closes the container after the command is completed.
func Execute(ctx context.Context, c *di.Container, cmd *cobra.Command) (err error) {
	ctx = di.ContextWithContainer(ctx, c)

	defer func() {
		err = multierr.Append(
			err,
			c.Close(),
		)
	}()

	return cmd.ExecuteContext(ctx)
}

// Setup returns a function that sets up the DI container to be used for a
// specific Cobra CLI command invocation. It should be assigned to the
// PersistentPreRunE field of the root command.
func Setup(fn interface{}) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		c := di.ContainerFromContext(ctx)

		c.Provide(func() (
			context.Context,
			*cobra.Command,
			[]string,
		) {
			return ctx, cmd, args
		})

		return c.Invoke(fn)
	}
}

// Run returns a function that invokes fn with arguments populated by the
// container. The returned function matches the signature of cobra.Command.Run.
func Run(fn interface{}) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if err := di.ContainerFromContext(ctx).Invoke(fn); err != nil {
			panic(err)
		}
	}
}

// RunE returns a function that invokes fn with arguments populated by the
// container. The returned function matches the signature of cobra.Command.RunE.
func RunE(fn interface{}) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return di.ContainerFromContext(ctx).Invoke(fn)
	}
}
