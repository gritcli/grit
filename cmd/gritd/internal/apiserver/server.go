package apiserver

import (
	"context"
	"errors"
	"sort"

	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/api"
)

// Server is an implementation of api.APIServer
type Server struct {
	Sources []source.Source
}

var _ api.APIServer = (*Server)(nil)

// ListSources lists the configured repository sources.
func (s *Server) ListSources(ctx context.Context, _ *api.ListSourcesRequest) (*api.ListSourcesResponse, error) {
	res := &api.ListSourcesResponse{}

	for _, s := range s.Sources {
		res.Sources = append(res.Sources, &api.Source{
			Name:        s.Name(),
			Description: s.Description(),
		})
	}

	sort.Slice(res.Sources, func(i, j int) bool {
		return res.Sources[i].Name < res.Sources[j].Name
	})

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
