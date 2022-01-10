package github

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/internal/daemon/internal/source"
	"github.com/gritcli/grit/internal/daemon/internal/source/internal/git"
)

// NewCloner returns a cloner that clones the repository with the given ID.
func (d *Driver) NewCloner(
	ctx context.Context,
	id string,
	clientLog logging.Logger,
) (source.Cloner, string, error) {
	serverLog := logging.Prefix(d.Logger, "clone[%s]: ", id)

	intID, err := parseRepoID(id)
	if err != nil {
		logging.LogString(serverLog, err.Error())
		return nil, "", err
	}

	r, ok := d.cache.RepoByID(intID)
	if !ok {
		var err error
		r, _, err = d.client.Repositories.GetByID(ctx, intID)
		if err != nil {
			logging.Log(serverLog, "unable to query API: %s", err)
			return nil, "", err
		}
	}

	c := &git.Cloner{
		Config:       d.Config.Git,
		SSHEndpoint:  r.GetSSHURL(),
		HTTPEndpoint: r.GetCloneURL(),
		Logger:       clientLog,
	}

	if d.Config.Token != "" {
		c.HTTPPassword = d.Config.Token
	}

	return c, r.GetFullName(), nil
}
