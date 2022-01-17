package githubsource

import (
	"context"

	"github.com/dogmatiq/dodeca/logging"
)

// Run performs any background processing required by the driver.
func (d *impl) Run(
	ctx context.Context,
	logger logging.Logger,
) error {
	return nil
}
