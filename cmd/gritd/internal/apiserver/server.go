package apiserver

import (
	"context"
	"errors"

	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/config"
)

// Server is an implementation of api.APIServer
type Server struct {
	Config config.Config
}

var _ api.APIServer = (*Server)(nil)

// ListSources lists the configured repository sources.
func (s *Server) ListSources(ctx context.Context, _ *api.ListSourcesRequest) (*api.ListSourcesResponse, error) {
	res := &api.ListSourcesResponse{}

	for _, s := range s.Config.Sources {
		res.Sources = append(res.Sources, &api.Source{
			Name:   s.Name,
			Driver: string(s.Config.Driver()),
			Config: s.Config.String(),
		})
	}

	return res, nil
}

// SearchRepositories looks for a repository by (partial) name.
func (s *Server) SearchRepositories(req *api.SearchRepositoriesRequest, stream api.API_SearchRepositoriesServer) error {
	return errors.New("not implemented")
}

// CloneRepository clones a remote repository.
func (s *Server) CloneRepository(ctx context.Context, req *api.CloneRepositoryRequest) (*api.CloneRepositoryResponse, error) {
	return nil, errors.New("not implemented")
}
