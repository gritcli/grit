package githubsource

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/builtins/gitvcs"
	"github.com/gritcli/grit/driver/sourcedriver"
)

// NewCloner returns a cloner that clones the repository with the given ID, and
// information about the repository being cloned.
func (s *source) NewCloner(
	ctx context.Context,
	id string,
	logger logging.Logger,
) (sourcedriver.Cloner, sourcedriver.RemoteRepo, error) {
	intID, err := parseRepoID(id)
	if err != nil {
		return nil, sourcedriver.RemoteRepo{}, err
	}

	r, ok := s.reposByID[intID]
	if !ok {
		var err error
		r, _, err = s.client.Repositories.GetByID(ctx, intID)
		if err != nil {
			return nil, sourcedriver.RemoteRepo{}, err
		}
	}

	logging.Debug(
		logger,
		"resolved %s to %s",
		id,
		r.GetFullName(),
	)

	c := &gitvcs.Cloner{
		SSHEndpoint:      r.GetSSHURL(),
		SSHKeyFile:       s.config.Git.SSHKeyFile,
		SSHKeyPassphrase: s.config.Git.SSHKeyPassphrase,
		HTTPEndpoint:     r.GetCloneURL(),
		PreferHTTP:       s.config.Git.PreferHTTP,
	}

	if s.config.Token != "" {
		c.HTTPPassword = s.config.Token
	}

	return c, toRemoteRepo(r), nil
}
