package apiserver

import (
	"context"
	"sort"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

// Server is the implementation of api.APIServer
type Server struct {
	SourceList source.List
	Cloner     *source.Cloner
	Logger     logging.Logger
}

// Sources lists the configured repository sources.
func (s *Server) Sources(ctx context.Context, _ *api.SourcesRequest) (*api.SourcesResponse, error) {
	res := &api.SourcesResponse{}

	for _, s := range s.SourceList {
		status, err := s.Driver.Status(ctx)
		if err != nil {
			return nil, err
		}

		res.Sources = append(res.Sources, &api.Source{
			Name:        s.Name,
			Description: s.Description,
			Status:      status,
			CloneDir:    s.CloneDir,
		})
	}

	sort.Slice(res.Sources, func(i, j int) bool {
		return res.Sources[i].Name < res.Sources[j].Name
	})

	return res, nil
}

// Resolve resolves repository name, URL or other identifier to a list of
// candidate repositories.
func (s *Server) Resolve(req *api.ResolveRequest, stream api.API_ResolveServer) error {
	ctx := stream.Context()
	g, ctx := errgroup.WithContext(ctx)

	log := s.newStreamLogger(
		stream,
		req.ClientOptions,
		func(out *api.ClientOutput) proto.Message {
			return &api.ResolveResponse{
				Response: &api.ResolveResponse_Output{
					Output: out,
				},
			}
		},
	)

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
				if err := stream.Send(&api.ResolveResponse{
					Response: &api.ResolveResponse_Repo{
						Repo: &api.Repo{
							Id:          r.ID,
							Source:      src.Name,
							Name:        r.Name,
							Description: r.Description,
							WebUrl:      r.WebURL,
						},
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

// Clone makes a local clone of a repository from a source.
func (s *Server) Clone(req *api.CloneRequest, stream api.API_CloneServer) error {
	dir, err := s.Cloner.Clone(
		stream.Context(),
		req.Source,
		req.RepoId,
		s.newStreamLogger(
			stream,
			req.ClientOptions,
			func(out *api.ClientOutput) proto.Message {
				return &api.CloneResponse{
					Response: &api.CloneResponse_Output{
						Output: out,
					},
				}
			},
		),
	)
	if err != nil {
		return err
	}

	return stream.Send(&api.CloneResponse{
		Response: &api.CloneResponse_Directory{
			Directory: dir,
		},
	})
}
