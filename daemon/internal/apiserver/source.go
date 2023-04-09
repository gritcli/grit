package apiserver

import (
	"context"
	"fmt"

	"github.com/gritcli/grit/api"
	"golang.org/x/exp/slices"
)

// ListSources lists the configured repository sources.
func (s *Server) ListSources(
	ctx context.Context,
	_ *api.ListSourcesRequest,
) (*api.ListSourcesResponse, error) {
	res := &api.ListSourcesResponse{}

	for _, src := range s.SourceList {
		status, err := src.Driver.Status(
			ctx,
			src.Log(s.Log),
		)
		if err != nil {
			return nil, err
		}

		res.Sources = append(res.Sources, &api.Source{
			Name:         src.Name,
			Description:  src.Description,
			Status:       status,
			BaseCloneDir: src.BaseCloneDir,
		})
	}

	slices.SortFunc(
		res.Sources,
		func(a, b *api.Source) bool {
			return a.Name < b.Name
		},
	)

	return res, nil
}

// SignIn signs in to a repository source.
func (s *Server) SignIn(
	req *api.SignInRequest,
	stream api.API_SignInServer,
) error {
	src, ok := s.SourceList.ByName(req.GetSource())
	if !ok {
		return fmt.Errorf("unable to sign in: unrecognized source (%s)", req.GetSource())
	}

	ctx := stream.Context()
	log := src.Log(s.Log)

	_, err := src.Driver.SignIn(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

// SignOut signs out of a repository source.
func (s *Server) SignOut(
	ctx context.Context,
	req *api.SignOutRequest,
) (*api.SignOutResponse, error) {
	src, ok := s.SourceList.ByName(req.GetSource())
	if !ok {
		return nil, fmt.Errorf("unable to sign out: unrecognized source (%s)", req.GetSource())
	}

	log := src.Log(s.Log)

	if err := src.Driver.SignOut(ctx, log); err != nil {
		return nil, err
	}

	return &api.SignOutResponse{}, nil
}
