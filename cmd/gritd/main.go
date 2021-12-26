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
		log logging.Logger,
	) error {
		go func() {
			<-ctx.Done()
			s.GracefulStop()
		}()

		lis, err := apiserver.Listen(cfg.Daemon.Socket)
		if err != nil {
			return err
		}
		defer lis.Close()

		logging.Log(log, "grit daemon v%s, listening for API requests at %s", version, cfg.Daemon.Socket)
		for _, src := range sources {
			logging.Log(log, "using '%s' repository source: %s", src.Name(), src.Description())
		}

		return s.Serve(lis)
	})
}
