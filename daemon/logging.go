package daemon

import (
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/imbue"
)

func init() {
	imbue.With0(
		container,
		func(
			ctx imbue.Context,
		) (logging.Logger, error) {
			return logging.DebugLogger, nil
		},
	)
}
