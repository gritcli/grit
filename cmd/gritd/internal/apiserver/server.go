package apiserver

import (
	"context"

	"github.com/gritcli/grit/internal/api"
)

// Server is a Grit gRPC API server.
type Server struct {
	Version string
}

// Ping is a no-op used to test that the server is reachable.
func (s *Server) Ping(ctx context.Context, _ *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{
		Version: s.Version,
	}, nil
}
