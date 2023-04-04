package apiserver

import (
	"context"

	"github.com/gritcli/grit/api"
	"golang.org/x/exp/maps"
)

// SuggestRepos returns a list of repository names to be used as suggestions for
// completing a partial repository name.
func (s *Server) SuggestRepos(
	ctx context.Context,
	req *api.SuggestReposRequest,
) (*api.SuggestResponse, error) {
	suggestions := s.Suggester.Suggest(
		req.Word,
		hasLocality(req.LocalityFilter, api.Locality_REMOTE),
		hasLocality(req.LocalityFilter, api.Locality_LOCAL),
	)

	return &api.SuggestResponse{
		Words: maps.Keys(suggestions),
	}, nil
}
