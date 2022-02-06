package apiserver

import (
	"context"
	"sort"

	"github.com/gritcli/grit/api"
)

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
