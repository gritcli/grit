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
	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/config"
	"golang.org/x/sync/errgroup"
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
		sources []source.Source,
		logger logging.Logger,
	) error {
		g, ctx := errgroup.WithContext(ctx)

		logging.Log(logger, "grit daemon v%s, listening for API requests at %s", version, cfg.Daemon.Socket)

		g.Go(func() error {
			lis, err := apiserver.Listen(cfg.Daemon.Socket)
			if err != nil {
				return err
			}
			defer lis.Close()

			go func() {
				<-ctx.Done()
				s.GracefulStop()
			}()

			return s.Serve(lis)
		})

		for _, src := range sources {
			src := src // capture loop variable

			g.Go(func() error {
				return src.Run(ctx)
			})
		}

		return g.Wait()
	})
}
