package apiserver

import "github.com/gritcli/grit/internal/api"

// Server is a Grit gRPC API server.
type Server struct{}

var _ api.APIServer = (*Server)(nil)
