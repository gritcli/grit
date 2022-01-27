package daemon

import "github.com/dogmatiq/dodeca/logging"

func init() {
	container.Provide(func() logging.Logger {
		return logging.DebugLogger
	})
}
