package apiserver

import (
	"context"

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
