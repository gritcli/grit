package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/cmd/gritd/internal/deps"
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

	return deps.Container.Invoke(func(
		cfg config.Config,
		s *grpc.Server,
		log logging.Logger,
	) error {
		go func() {
			<-ctx.Done()
			s.GracefulStop()
		}()

		logging.Log(log, "grit daemon v%s, listening for API requests at %s", version, cfg.Daemon.Socket)
		lis, err := apiserver.Listen(cfg.Daemon.Socket)
		if err != nil {
			return err
		}
		defer lis.Close()

		return s.Serve(lis)
	})
}
