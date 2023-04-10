package githubsource

import (
	"context"
	"errors"

	"github.com/gritcli/grit/daemon/internal/logs"
)

// SignIn signs in to the source.
func (s *source) SignIn(
	ctx context.Context,
	log logs.Log,
) error {
	if s.config.Token != "" {
		return errors.New("already authenticated using a personal access token (PAT)")
	}
	return errors.New("<not implemented>")
}

// SignOut signs out of the source.
func (s *source) SignOut(
	ctx context.Context,
	log logs.Log,
) error {
	return errors.New("<not implemented>")
}
