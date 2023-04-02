package githubsource

import (
	"context"

	"github.com/gritcli/grit/logs"
)

// Run performs any background processing required by the source.
func (s *source) Run(
	ctx context.Context,
	log logs.Log,
) error {
	return nil
}
