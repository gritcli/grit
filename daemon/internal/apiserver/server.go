package apiserver

import (
	"errors"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/logs"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// Server is the implementation of api.APIServer
type Server struct {
	SourceList source.List
	Cloner     *source.Cloner
	Suggester  *source.Suggester
	Log        logs.Log
}

// newLogger returns a logs.Logger that sends messages over a gRPC stream.
func (s *Server) newLogger(
	stream grpc.ServerStream,
	options *api.ClientOptions,
	wrap func(*api.ClientOutput) proto.Message,
) logs.Log {
	return func(m logs.Message) {
		if m.IsVerbose && !options.GetVerbose() {
			return
		}

		msg := wrap(
			&api.ClientOutput{
				Message: m.Text,
			},
		)

		if err := stream.SendMsg(msg); err != nil {
			s.Log.Write("unable to write message to stream: %s", err)
		}
	}
}

// Listen starts a listener on the given unix socket.
//
// It deletes the socket file if it already exists.
func Listen(socket string) (net.Listener, error) {
	l, err := net.Listen("unix", socket)
	if err == nil {
		return l, nil
	}

	if !errors.Is(err, syscall.EADDRINUSE) {
		return nil, err
	}

	if err := os.Remove(socket); err != nil {
		return nil, err
	}

	return net.Listen("unix", socket)
}

// marshalRemoteRepo marshals a sourcedriver.RemoteRepo into its API
// representation.
func marshalRemoteRepo(source string, r sourcedriver.RemoteRepo) *api.RemoteRepo {
	return &api.RemoteRepo{
		Id:          r.ID,
		Source:      source,
		Name:        r.Name,
		Description: r.Description,
		WebUrl:      r.WebURL,
	}
}

// marshalRemoteRepo marshals a source.LocalRepo into its API
// representation.
func marshalLocalRepo(r source.LocalRepo) *api.LocalRepo {
	return &api.LocalRepo{
		RemoteRepo:       marshalRemoteRepo(r.Source.Name, r.RemoteRepo),
		AbsoluteCloneDir: r.AbsoluteCloneDir,
	}
}

// hasLocality returns true if the filter includes the given locality.
//
// An empty filter is considered to contain all localities.
func hasLocality(filter []api.Locality, loc api.Locality) bool {
	if len(filter) == 0 {
		return true
	}

	for _, l := range filter {
		if l == loc {
			return true
		}
	}

	return false
}

// hasSource returns true if the filter includes the given source.
//
// An empty filter is considered to contain all sources.
func hasSource(filter []string, source string) bool {
	if len(filter) == 0 {
		return true
	}

	for _, s := range filter {
		if strings.EqualFold(s, source) {
			return true
		}
	}

	return false
}
