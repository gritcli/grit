package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/deps"
	"github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/driver/registry"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Run executes the Grit daemon.
func Run(version string) (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	return deps.Container.Invoke(func(
		cfg config.Config,
		s *grpc.Server,
		r *registry.Registry,
		sources source.List,
		logger logging.Logger,
	) error {
		logging.Log(logger, "grit daemon v%s", version)

		logDrivers(logger, r)
		logSources(logger, sources)

		if err := initSourceDrivers(ctx, logger, sources); err != nil {
			return err
		}

		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return runSourceDrivers(ctx, logger, sources)
		})

		g.Go(func() error {
			return runGRPCServer(ctx, logger, cfg.Daemon.Socket, s)
		})

		return g.Wait()
	})
}

// logDrivers logs information about the drivers in the registry.
func logDrivers(logger logging.Logger, r *registry.Registry) {
	for alias, reg := range r.SourceDrivers() {
		if alias == reg.Name {
			logger.Log(
				"config: loaded '%s' source driver: %s",
				reg.Name,
				reg.Description,
			)
		} else {
			logger.Log(
				"config: loaded '%s' source driver as '%s': %s",
				reg.Name,
				reg.Description,
				alias,
			)
		}
	}

	for alias, reg := range r.VCSDrivers() {
		if alias == reg.Name {
			logger.Log(
				"config: loaded '%s' vcs driver: %s",
				reg.Name,
				reg.Description,
			)
		} else {
			logger.Log(
				"config: loaded '%s' vcs driver as '%s': %s",
				reg.Name,
				reg.Description,
				alias,
			)
		}
	}
}

// logSources logs information about the sources in the configuration.
func logSources(logger logging.Logger, sources source.List) {
	for _, s := range sources {
		logger.Log(
			"config: loaded '%s' source: %s",
			s.Name,
			s.Description,
		)
	}
}

// initSourceDrivers initializes each source's driver in parallel.
func initSourceDrivers(ctx context.Context, logger logging.Logger, sources source.List) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range sources {
		src := src // capture loop variable
		g.Go(func() error {
			return src.Driver.Init(
				ctx,
				logging.Prefix(
					logger,
					"source[%s]: ",
					src.Name,
				),
			)
		})
	}

	return g.Wait()
}

// runSourceDrivers runs each source's driver in parallel.
func runSourceDrivers(ctx context.Context, logger logging.Logger, sources source.List) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range sources {
		src := src // capture loop variable
		g.Go(func() error {
			return src.Driver.Run(
				ctx,
				logging.Prefix(
					logger,
					"source[%s]: ",
					src.Name,
				),
			)
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
