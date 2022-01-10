package apiserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

// server is an implementation of api.APIServer
type server struct {
	sources source.List
}

// New returns a new API server.
func New(sources source.List) api.APIServer {
	return &server{sources}
}

// Sources lists the configured repository sources.
func (s *server) Sources(ctx context.Context, _ *api.SourcesRequest) (*api.SourcesResponse, error) {
	res := &api.SourcesResponse{}

	for _, s := range s.sources {
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
func (s *server) Resolve(req *api.ResolveRequest, stream api.API_ResolveServer) error {
	ctx := stream.Context()
	g, ctx := errgroup.WithContext(ctx)

	log := newStreamLogger(
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

	for _, src := range s.sources {
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
func (s *server) Clone(req *api.CloneRequest, stream api.API_CloneServer) error {
	ctx := stream.Context()

	src, ok := s.sources.ByName(req.Source)
	if !ok {
		return errors.New("unrecognized source name")
	}

	cloner, dir, err := src.Driver.NewCloner(
		ctx,
		req.RepoId,
		newStreamLogger(
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

	dir = filepath.Join(src.CloneDir, dir)

	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("clone directory (%s) already exists", dir)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := cloner.Clone(ctx, dir); err != nil {
		os.RemoveAll(dir)
		return err
	}

	return stream.Send(&api.CloneResponse{
		Response: &api.CloneResponse_Directory{
			Directory: dir,
		},
	})
}
