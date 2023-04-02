package daemon

import (
	"fmt"
	"os"

	"github.com/dogmatiq/imbue"
	"github.com/gritcli/grit/logs"
)

func init() {
	imbue.With0(
		container,
		func(
			ctx imbue.Context,
		) (logs.Log, error) {
			return func(m logs.Message) {
				fmt.Fprintln(os.Stdout, m.Text)
			}, nil
		},
	)
}
