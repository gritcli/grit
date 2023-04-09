package apiserver

import (
	"context"

	"github.com/gritcli/grit/api"
)

// DaemonInfo returns information about the daemon.
func (s *Server) DaemonInfo(
	ctx context.Context,
	req *api.DaemonInfoRequest,
) (*api.DaemonInfoResponse, error) {
	return &api.DaemonInfoResponse{
		Version: s.Version,
		Pid:     uint64(s.PID),
	}, nil
}
