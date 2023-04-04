package daemonapiserver

import (
	"context"

	"github.com/gritcli/grit/api/daemonapi"
)

// Server is an implementation of daemonapi.APIServer.
type Server struct {
	Version string
	PID     int
}

// Info returns information about the daemon.
func (s *Server) Info(
	ctx context.Context,
	req *daemonapi.InfoRequest,
) (*daemonapi.InfoResponse, error) {
	return &daemonapi.InfoResponse{
		Version: s.Version,
		Pid:     uint64(s.PID),
	}, nil
}
