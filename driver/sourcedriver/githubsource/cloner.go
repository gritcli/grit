package githubsource

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/driver/sourcedriver"
	"github.com/gritcli/grit/driver/vcsdriver/gitvcs"
)

// NewCloner returns a cloner that clones the repository with the given ID.
func (d *impl) NewCloner(
	ctx context.Context,
	id string,
	logger logging.Logger,
) (sourcedriver.Cloner, string, error) {
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

	c := &gitvcs.Cloner{
		SSHEndpoint:      r.GetSSHURL(),
		SSHKeyFile:       d.config.Git.SSHKeyFile,
		SSHKeyPassphrase: d.config.Git.SSHKeyPassphrase,
		HTTPEndpoint:     r.GetCloneURL(),
		PreferHTTP:       d.config.Git.PreferHTTP,
	}

	if d.config.Token != "" {
		c.HTTPPassword = d.config.Token
	}

	return c, r.GetFullName(), nil
}
