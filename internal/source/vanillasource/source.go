package vanillasource

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport"
)

// Source is an implementation of source.Source for "vanilla" Git servers.
type Source struct {
	Endpoint *transport.Endpoint
}

// Description returns a short, human-readable description of the source.
//
// The description should be adequate to distinguish this source from any
// other sources that may exist.
func (s *Source) Description() string {
	return fmt.Sprintf(
		"git %s",
		s.Endpoint.Host,
	)
}

// Status queries the status of the source.
//
// It returns an error if the source is misconfigured or unreachable.
//
// The status string should include any source-specific information
func (s *Source) Status(ctx context.Context) (string, error) {
	return s.Endpoint.Host, errors.New("not implemented")
}
