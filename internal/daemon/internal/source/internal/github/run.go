package github

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
)

// Run performs any ongoing behavior required by the driver.
func (d *Driver) Run(
	ctx context.Context,
	logger logging.Logger,
) error {
	return nil
}
