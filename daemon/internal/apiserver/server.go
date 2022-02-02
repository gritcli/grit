package apiserver

import (
	"context"
	"sort"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/api"
	"github.com/gritcli/grit/daemon/internal/source"
	"github.com/gritcli/grit/driver/sourcedriver"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

// Server is the implementation of api.APIServer
type Server struct {
	SourceList source.List
	Cloner     *source.Cloner
	Suggester  *source.Suggester
	Logger     logging.Logger
}

// ListSources lists the configured repository sources.
func (s *Server) ListSources(ctx context.Context, _ *api.ListSourcesRequest) (*api.ListSourcesResponse, error) {
	res := &api.ListSourcesResponse{}

	for _, s := range s.SourceList {
		status, err := s.Driver.Status(ctx)
		if err != nil {
			return nil, err
		}

		res.Sources = append(res.Sources, &api.Source{
			Name:         s.Name,
			Description:  s.Description,
			Status:       status,
			BaseCloneDir: s.BaseCloneDir,
		})
	}

	sort.Slice(res.Sources, func(i, j int) bool {
		return res.Sources[i].Name < res.Sources[j].Name
	})

	return res, nil
}

// ResolveRepo resolves a repository name, URL or other identifier to a list of
// repositories.
func (s *Server) ResolveRepo(
	req *api.ResolveRepoRequest,
	stream api.API_ResolveRepoServer,
) error {
	ctx := stream.Context()
	g, ctx := errgroup.WithContext(ctx)

	log := s.newStreamLogger(
		stream,
		req.ClientOptions,
		func(out *api.ClientOutput) proto.Message {
			return &api.ResolveRepoResponse{
				Response: &api.ResolveRepoResponse_Output{
					Output: out,
				},
			}
		},
	)

	if req.Locality != api.Locality_LOCAL_ONLY {
		for _, src := range s.SourceList {
			src := src // capture loop variable

			g.Go(func() error {
				repos, err := src.Driver.Resolve(
					ctx,
					req.Query,
					logging.Prefix(log, "%s: ", src.Name),
				)
				if err != nil {
					return err
				}

				for _, r := range repos {
					if err := stream.Send(&api.ResolveRepoResponse{
						Response: &api.ResolveRepoResponse_RemoteRepo{
							RemoteRepo: marshalRemoteRepo(src.Name, r),
						},
					}); err != nil {
						return err
					}
				}

				return nil
			})
		}
	}

	return g.Wait()
}

// CloneRepo makes a local clone of a repository from a source.
func (s *Server) CloneRepo(req *api.CloneRepoRequest, stream api.API_CloneRepoServer) error {
	repo, err := s.Cloner.Clone(
		stream.Context(),
		req.Source,
		req.RepoId,
		s.newStreamLogger(
			stream,
			req.ClientOptions,
			func(out *api.ClientOutput) proto.Message {
				return &api.CloneRepoResponse{
					Response: &api.CloneRepoResponse_Output{
						Output: out,
					},
				}
			},
		),
	)
	if err != nil {
		return err
	}

	return stream.Send(&api.CloneRepoResponse{
		Response: &api.CloneRepoResponse_LocalRepo{
			LocalRepo: marshalLocalRepo(repo),
		},
	})
}

// SuggestRepos returns a list of repository names to be used as suggestions for
// completing a partial repository name.
func (s *Server) SuggestRepos(
	ctx context.Context,
	req *api.SuggestReposRequest,
) (*api.SuggestResponse, error) {
	repos := s.Suggester.Suggest(
		req.Word,
		req.Locality != api.Locality_REMOTE_ONLY,
		req.Locality != api.Locality_LOCAL_ONLY,
	)

	res := &api.SuggestResponse{}
	for _, r := range repos {
		res.Words = append(res.Words, r.Name)
	}

	return res, nil
}

// marshalRemoteRepo marshals a sourcedriver.RemoteRepo into its API
// representation.
func marshalRemoteRepo(source string, r sourcedriver.RemoteRepo) *api.RemoteRepo {
	return &api.RemoteRepo{
		Id:          r.ID,
		Source:      source,
		Name:        r.Name,
		Description: r.Description,
		WebUrl:      r.WebURL,
	}
}

// marshalRemoteRepo marshals a source.LocalRepo into its API
// representation.
func marshalLocalRepo(r source.LocalRepo) *api.LocalRepo {
	return &api.LocalRepo{
		RemoteRepo:       marshalRemoteRepo(r.Source.Name, r.RemoteRepo),
		AbsoluteCloneDir: r.AbsoluteCloneDir,
	}
}
