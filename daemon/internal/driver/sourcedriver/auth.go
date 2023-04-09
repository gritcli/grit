package sourcedriver

import (
	"context"

	"github.com/gritcli/grit/daemon/internal/logs"
)

// Authenticator is an interface for authenticating a user with a source.
type Authenticator interface {
	Authenticate(
		ctx context.Context,
		log logs.Log,
	) error
}
