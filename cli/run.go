package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/commands"
)

// container is the dependency injection container for the Grit CLI.
var container = imbue.New()

// Run starts the Grit CLI.
func Run(version string) (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	return commands.
		Root(container, version).
		ExecuteContext(ctx)
}
