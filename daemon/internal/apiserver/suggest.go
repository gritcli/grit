package apiserver

import (
	"context"

	"github.com/gritcli/grit/api"
)

// SuggestRepos returns a list of repository names to be used as suggestions for
// completing a partial repository name.
func (s *Server) SuggestRepos(
	ctx context.Context,
	req *api.SuggestReposRequest,
) (*api.SuggestResponse, error) {
	repos := s.Suggester.Suggest(
		req.Word,
		hasLocality(req.LocalityFilter, api.Locality_REMOTE),
		hasLocality(req.LocalityFilter, api.Locality_LOCAL),
	)

	res := &api.SuggestResponse{}
	for _, r := range repos {
		res.Words = append(res.Words, r.Name)
	}

	return res, nil
}
