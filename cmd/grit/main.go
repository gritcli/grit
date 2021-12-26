package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gritcli/grit/cmd/grit/internal/commands"
	"github.com/gritcli/grit/cmd/grit/internal/deps"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := run(); err != nil {
		os.Exit(1)
	}
}

// version string, automatically set during build process.
var version = "0.0.0"

func run() (err error) {
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
