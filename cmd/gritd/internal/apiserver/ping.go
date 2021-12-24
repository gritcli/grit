package apiserver

import (
	"context"

	"github.com/gritcli/grit/internal/api"
)

// PingServer is an implementation of api.PingServer
type PingServer struct {
	Version string
}

// Ping is a no-op used to test that the server is reachable.
func (s *PingServer) Ping(ctx context.Context, _ *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{
		Version: s.Version,
	}, nil
}
