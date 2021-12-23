package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/di"
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

	container := di.New()

	container.Provide(func() *apiserver.Server {
		return &apiserver.Server{}
	})

	return container.Invoke(func(s *apiserver.Server) error {
		g := grpc.NewServer()
		api.RegisterAPIServer(g, s)

		go func() {
			<-ctx.Done()
			g.GracefulStop()
		}()

		l, err := net.Listen("tcp", "/var/run/gritd.socket")
		if err != nil {
			return err
		}
		defer l.Close()

		return g.Serve(l)
	})
}
