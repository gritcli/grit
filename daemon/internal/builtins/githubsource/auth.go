package githubsource

import (
	"context"
	"errors"

	"github.com/gritcli/grit/daemon/internal/driver/sourcedriver"
	"github.com/gritcli/grit/daemon/internal/logs"
)

// SignIn signs in to the source.
func (s *source) SignIn(
	ctx context.Context,
	log logs.Log,
) (sourcedriver.Authenticator, error) {
	if s.config.Token != "" {
		return nil, errors.New("already authenticated using a personal access token (PAT)")
	}

	return nil, errors.New("<not implemented>")
}

// SignOut signs out of the source.
func (s *source) SignOut(
	ctx context.Context,
	log logs.Log,
) error {
	return errors.New("<not implemented>")
}
