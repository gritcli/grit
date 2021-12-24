package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogmatiq/dapper"
	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/cmd/gritd/internal/deps"
	"github.com/gritcli/grit/internal/commondeps"
	"github.com/gritcli/grit/internal/config"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
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

	commondeps.Provide(&deps.Container, version)

	return deps.Container.Invoke(func(
		cfg config.Config,
		s *grpc.Server,
	) error {
		dapper.Print(cfg)

		go func() {
			<-ctx.Done()
			s.GracefulStop()
		}()

		l, err := apiserver.Listen(cfg.Daemon.Socket)
		if err != nil {
			return err
		}
		defer l.Close()

		return s.Serve(l)
	})
}
