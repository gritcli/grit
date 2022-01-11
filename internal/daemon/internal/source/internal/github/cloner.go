package github

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/daemon/internal/source/internal/git"
)

// NewBoundCloner returns a bound cloner that clones the repository with the
// given ID.
func (d *Driver) NewBoundCloner(
	ctx context.Context,
	id string,
	logger logging.Logger,
) (source.BoundCloner, string, error) {
	intID, err := parseRepoID(id)
	if err != nil {
		return nil, "", err
	}

	r, ok := d.cache.RepoByID(intID)
	if !ok {
		var err error
		r, _, err = d.client.Repositories.GetByID(ctx, intID)
		if err != nil {
			return nil, "", err
		}
	}

	logging.Debug(
		logger,
		"resolved %s to %s",
		id,
		r.GetFullName(),
	)

	c := &git.BoundCloner{
		Config:       d.Config.Git,
		SSHEndpoint:  r.GetSSHURL(),
		HTTPEndpoint: r.GetCloneURL(),
	}

	if d.Config.Token != "" {
		c.HTTPPassword = d.Config.Token
	}

	return c, r.GetFullName(), nil
}
