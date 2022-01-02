package apiserver

import (
	"context"
	"sort"

	"github.com/gritcli/grit/cmd/gritd/internal/source"
	"github.com/gritcli/grit/internal/api"
	"golang.org/x/sync/errgroup"
)

// server is an implementation of api.APIServer
type server struct {
	sources []source.Source
}

// New returns a new API server.
func New(sources []source.Source) api.APIServer {
	return &server{sources}
}

// Sources lists the configured repository sources.
func (s *server) Sources(ctx context.Context, _ *api.SourcesRequest) (*api.SourcesResponse, error) {
	res := &api.SourcesResponse{}

	for _, s := range s.sources {
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

// Resolve resolves repository name, URL or other identifier to a list of
// candidate repositories.
func (s *server) Resolve(req *api.ResolveRequest, stream api.API_ResolveServer) error {
	ctx := stream.Context()
	g, ctx := errgroup.WithContext(ctx)

	for _, src := range s.sources {
		src := src // capture loop variable

		g.Go(func() error {
			repos, err := src.Resolve(ctx, req.Query)
			if err != nil {
				return err
			}

			for _, r := range repos {
				if err := stream.Send(&api.ResolveResponse{
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
