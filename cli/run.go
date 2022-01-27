package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gritcli/grit/cli/internal/cobradi"
	"github.com/gritcli/grit/cli/internal/commands"
	"github.com/gritcli/grit/internal/di"
)

// container is the dependency injection container for the Grit CLI.
var container di.Container

// Run starts the Grit CLI.
func Run(version string) (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	return cobradi.Execute(
		ctx,
		&container,
		commands.NewRoot(version),
	)
}
