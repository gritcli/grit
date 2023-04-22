package githubsource

import (
	"context"

	"github.com/gritcli/grit/daemon/internal/logs"
)

// Run performs any background processing required by the source.
func (s *source) Run(
	ctx context.Context,
	log logs.Log,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.reinit:
			if err := s.init(ctx); err != nil {
				return err
			}
		}
	}
}
