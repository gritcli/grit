package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/source"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// container is the dependency injection container for the Grit daemon.
var container = imbue.New()

// Run executes the Grit daemon.
func Run(version string) (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	if err := imbue.Invoke3(
		ctx,
		container,
		func(
			ctx context.Context,
			r *config.DriverRegistry,
			s source.List,
			l logging.Logger,
		) error {
			logging.Log(l, "grit daemon v%s", version)

			logDrivers(l, r)
			logSources(l, s)

			return initSourceDrivers(ctx, l, s)
		},
	); err != nil {
		return err
	}

	g := container.WaitGroup(ctx)
	imbue.Go2(g, runSourceDrivers)
	imbue.Go3(g, runGRPCServer)

	return g.Wait()
}

// logDrivers logs information about the drivers in the registry.
func logDrivers(logger logging.Logger, r *config.DriverRegistry) {
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
func runSourceDrivers(
	ctx context.Context,
	logger logging.Logger,
	sources source.List,
) error {
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
	cfg config.Config,
	s *grpc.Server,
	l logging.Logger,
) error {
	lis, err := apiserver.Listen(cfg.Daemon.Socket)
	if err != nil {
		return err
	}
	defer lis.Close()

	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()

	logging.Log(l, "api: accepting requests on unix socket: %s", cfg.Daemon.Socket)
	defer logging.Log(l, "api: server stopped")

	return s.Serve(lis)
}
