package apiserver

import (
	"context"

	"github.com/gritcli/grit/internal/api"
	"github.com/gritcli/grit/internal/config"
)

// SourceAPIServer is an implementation of api.SourceAPIServer
type SourceAPIServer struct {
	Config config.Config
}

// ListSources lists the configured repository sources.
func (s *SourceAPIServer) ListSources(ctx context.Context, _ *api.ListSourcesRequest) (*api.ListSourcesResponse, error) {
	res := &api.ListSourcesResponse{}

	for _, s := range s.Config.Sources {
		res.Sources = append(res.Sources, &api.Source{
			Name:   s.Name,
			Driver: string(s.Driver),
			Config: s.Config.DescribeConfig(),
		})
	}

	return res, nil
}
