package daemon

import (
	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/logs"
)

func init() {
	imbue.With0(
		catalog,
		func(
			ctx imbue.Context,
		) (logs.Log, error) {
			return logs.Verbose, nil
		},
	)
}
