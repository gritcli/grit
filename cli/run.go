package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/cli/internal/commands"
)

// catalog is the dependency injection catalog for the Grit CLI.
var catalog = imbue.NewCatalog()

// Run starts the Grit CLI.
func Run(version string) (err error) {
	con := imbue.New(imbue.WithCatalog(catalog))
	defer con.Close()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	return commands.
		Root(con, version).
		ExecuteContext(ctx)
}
