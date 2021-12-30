package apiserver

import (
	"context"
	"errors"
	"sort"

	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/api"
	"golang.org/x/sync/errgroup"
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

// ResolveRepoName resolves a repository name to a list of candidate
// repositories.
func (s *Server) ResolveRepoName(req *api.ResolveRepoNameRequest, stream api.API_ResolveRepoNameServer) error {
	ctx := stream.Context()
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range s.Sources {
		src := src // capture loop variable

		g.Go(func() error {
			repos, err := src.Resolve(ctx, req.Name)
			if err != nil {
				return err
			}

			for _, r := range repos {
				if err := stream.Send(&api.ResolveRepoNameResponse{
					Repo: &api.Repo{
						SourceName:  src.Name(),
						RepoId:      r.ID,
						RepoName:    r.Name,
						Description: r.Description,
						WebUrl:      r.WebURL,
					},
				}); err != nil {
					return err
				}
			}

			return nil
		})
	}

	return g.Wait()
}

// CloneRepository clones a remote repository.
func (s *Server) CloneRepository(ctx context.Context, req *api.CloneRepositoryRequest) (*api.CloneRepositoryResponse, error) {
	return nil, errors.New("not implemented")
}