package daemon

import (
	"context"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/daemon/internal/apiserver"
	"github.com/gritcli/grit/daemon/internal/config"
	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/httpserver"
	"github.com/gritcli/grit/daemon/internal/logs"
	"github.com/gritcli/grit/daemon/internal/signalx"
	"github.com/gritcli/grit/daemon/internal/source"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// catalog is the dependency injection catalog for the Grit daemon.
var catalog = imbue.NewCatalog()

// Run executes the Grit daemon.
func Run(ver string) error {
	imbue.With0Named[version](
		catalog,
		func(
			ctx imbue.Context,
		) (string, error) {
			return ver, nil
		},
	)

	reloads := 0
	for {
		reload, err := run(reloads)
		if !reload || err != nil {
			return err
		}

		reloads++
	}
}

// run executes the Grit daemon.
func run(reloads int) (reload bool, err error) {
	con := imbue.New(imbue.WithCatalog(catalog))
	defer con.Close()

	ctx, cancel := signalx.NotifyContextWithCause(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer func() {
		if ctx.Err() == context.Canceled {
			err = nil
			reload = signalx.SignalCause(ctx) == syscall.SIGHUP
		}
		cancel()
	}()

	if err := imbue.Invoke5(
		ctx,
		con,
		func(
			ctx context.Context,
			ver imbue.ByName[version, string],
			r *config.DriverRegistry,
			s source.List,
			lis imbue.ByName[httpListener, net.Listener],
			log logs.Log,
		) error {
			if reloads == 0 {
				log.Write("grit daemon v%s, pid %d", ver.Value(), os.Getpid())
			} else {
				log.Write("reloading daemon configuration, pid %d", os.Getpid())
			}

			logDrivers(r, log)
			logSources(s, log)

			return initSourceDrivers(ctx, s, lis.Value(), log)
		},
	); err != nil {
		return false, err
	}

	g := con.WaitGroup(ctx)
	imbue.Go2(g, runSourceDrivers)
	imbue.Go3(g, runGRPCServer)
	imbue.Go3(g, runHTTPServer)

	return false, g.Wait()
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
	lis net.Listener,
	log logs.Log,
) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range sources {
		src := src // capture loop variable

		g.Go(func() error {
			return src.Driver.Init(
				ctx,
				sourcedriver.InitParameters{
					BaseURL: src.BaseURL,
				},
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

	log.Write("api: accepting requests on %s", cfg.Daemon.Socket)
	defer log.Write("api: server stopped")

	if err := s.Serve(lis); err != nil {
		return err
	}

	return ctx.Err()
}

// runHTTPServer runs the HTTP server.
func runHTTPServer(
	ctx context.Context,
	lis imbue.ByName[httpListener, net.Listener],
	sources source.List,
	log logs.Log,
) error {
	mux := http.NewServeMux()
	mux.Handle("/", &httpserver.IndexHandler{})

	for _, s := range sources {
		mux.Handle(s.BaseURL.Path, s.Driver)
	}

	s := &http.Server{
		Handler:           mux,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		s.Close()
	}()

	listener := lis.Value()
	log.Write("http: accepting requests at http://%s", listener.Addr())
	defer log.Write("http: server stopped")

	if err := s.Serve(listener); err != nil {
		return err
	}

	return ctx.Err()
}
