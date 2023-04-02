package apiserver

import (
	"context"

	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/logs"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

// ResolveRepo resolves a repository name, URL or other identifier to a list of
// repositories.
func (s *Server) ResolveRepo(
	req *api.ResolveRepoRequest,
	responses api.API_ResolveRepoServer,
) error {
	ctx := responses.Context()
	g, ctx := errgroup.WithContext(ctx)

	logger := s.newLogger(
		responses,
		req.ClientOptions,
		func(out *api.ClientOutput) proto.Message {
			return &api.ResolveRepoResponse{
				Response: &api.ResolveRepoResponse_Output{
					Output: out,
				},
			}
		},
	)

	for _, src := range s.SourceList {
		src := src // capture loop variable

		if !hasSource(req.SourceFilter, src.Name) {
			continue
		}

		if hasLocality(req.LocalityFilter, api.Locality_REMOTE) {
			g.Go(func() error {
				return s.resolveRemoteRepo(
					ctx,
					src,
					req.Query,
					responses,
					logger.WithPrefix("%s: ", src.Name),
				)
			})
		}
	}

	return g.Wait()
}

// resolveRemoteRepo sends a response for each repository from src that matches
// the given query.
func (s *Server) resolveRemoteRepo(
	ctx context.Context,
	src source.Source,
	query string,
	responses api.API_ResolveRepoServer,
	log logs.Log,
) error {
	repos, err := src.Driver.Resolve(ctx, query, log)
	if err != nil {
		return err
	}

	for _, r := range repos {
		if err := responses.Send(&api.ResolveRepoResponse{
			Response: &api.ResolveRepoResponse_RemoteRepo{
				RemoteRepo: marshalRemoteRepo(src.Name, r),
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
