package apiserver

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/common/api"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
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

	src, ok := s.sourceByName(req.Source)
	if !ok {
		return errors.New("unrecognized source name")
	}

	tempDir, err := os.MkdirTemp("", "grit-clone-")
	if err != nil {
		return err
	}

	relDir, err := src.Driver.Clone(
		ctx,
		req.RepoId,
		tempDir,
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

	cloneDir := filepath.Join(src.CloneDir, relDir)
	parentDir := filepath.Dir(cloneDir)

	if err := os.MkdirAll(parentDir, 0700); err != nil {
		return err
	}

	if err := os.Rename(tempDir, cloneDir); err != nil {
		// TODO: check if error is due to tempDir and cloneDir being on
		// different disk drives and if so, copy, then delete.
		return err
	}

	return stream.Send(&api.CloneResponse{
		Response: &api.CloneResponse_Directory{
			Directory: cloneDir,
		},
	})
}

func (s *server) sourceByName(n string) (source.Source, bool) {
	for _, src := range s.sources {
		if strings.EqualFold(src.Name, n) {
			return src, true
		}
	}

	return source.Source{}, false
}
