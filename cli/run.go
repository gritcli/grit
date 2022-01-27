package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gritcli/grit/cli/internal/commands"
	"github.com/gritcli/grit/cli/internal/deps"
)

// Run starts the Grit CLI.
func Run(version string) (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	return deps.Execute(
		ctx,
		&deps.Container,
		commands.NewRoot(version),
	)
}
