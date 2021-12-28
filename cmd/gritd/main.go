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
		logging.Log(logger, "grit daemon v%s", version)

		if err := initSources(ctx, logger, sources); err != nil {
			return err
		}

		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return runSources(ctx, logger, sources)
		})

		g.Go(func() error {
			return runGRPCServer(ctx, logger, cfg.Daemon.Socket, s)
		})

		return g.Wait()
	})
}

// initSources initializes each source in parallel.
func initSources(ctx context.Context, logger logging.Logger, sources []source.Source) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range sources {
		src := src // capture loop variable
		g.Go(func() error {
			return src.Init(ctx)
		})
	}

	return g.Wait()
}

// runSources runs each source in parallel.
func runSources(ctx context.Context, logger logging.Logger, sources []source.Source) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range sources {
		src := src // capture loop variable
		g.Go(func() error {
			return src.Run(ctx)
		})
	}

	return g.Wait()
}

// runGRPCServer runs the gRPC server.
func runGRPCServer(
	ctx context.Context,
	logger logging.Logger,
	socket string,
	s *grpc.Server,
) error {
	lis, err := apiserver.Listen(socket)
	if err != nil {
		return err
	}
	defer lis.Close()

	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()

	logging.Log(logger, "api: accepting requests on unix socket: %s", socket)
	defer logging.Log(logger, "api: server stopped")

	return s.Serve(lis)
}
