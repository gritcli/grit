package deps

import "github.com/dogmatiq/dodeca/logging"

func init() {
	Container.Provide(func() logging.Logger {
		return logging.DefaultLogger
	})
}
