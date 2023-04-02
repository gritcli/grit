package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/config"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/logs"
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
			log logs.Log,
		) error {
			log.Write("grit daemon v%s", version)

			logDrivers(r, log)
			logSources(s, log)

			return initSourceDrivers(ctx, s, log)
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
func logDrivers(
	r *config.DriverRegistry,
	log logs.Log,
) {
	for alias, reg := range r.SourceDrivers() {
		if alias == reg.Name {
			log.Write(
				"config: loaded '%s' source driver: %s",
				reg.Name,
				reg.Description,
			)
		} else {
			log.Write(
				"config: loaded '%s' source driver as '%s': %s",
				reg.Name,
				reg.Description,
				alias,
			)
		}
	}

	for alias, reg := range r.VCSDrivers() {
		if alias == reg.Name {
			log.Write(
				"config: loaded '%s' vcs driver: %s",
				reg.Name,
				reg.Description,
			)
		} else {
			log.Write(
				"config: loaded '%s' vcs driver as '%s': %s",
				reg.Name,
				reg.Description,
				alias,
			)
		}
	}
}

// logSources logs information about the sources in the configuration.
func logSources(
	sources source.List,
	log logs.Log,
) {
	for _, s := range sources {
		log.Write(
			"config: loaded '%s' source: %s",
			s.Name,
			s.Description,
		)
	}
}

// initSourceDrivers initializes each source's driver in parallel.
func initSourceDrivers(
	ctx context.Context,
	sources source.List,
	log logs.Log,
) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range sources {
		src := src // capture loop variable
		g.Go(func() error {
			return src.Driver.Init(
				ctx,
				src.Log(log),
			)
		})
	}

	return g.Wait()
}

// runSourceDrivers runs each source's driver in parallel.
func runSourceDrivers(
	ctx context.Context,
	sources source.List,
	log logs.Log,
) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range sources {
		src := src // capture loop variable
		g.Go(func() error {
			return src.Driver.Run(
				ctx,
				src.Log(log),
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
	log logs.Log,
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

	log.Write("api: accepting requests on unix socket: %s", cfg.Daemon.Socket)
	defer log.Write("api: server stopped")

	return s.Serve(lis)
}
