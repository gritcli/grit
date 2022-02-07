package apitea

import "github.com/gritcli/grit/api"

// OutputReceived is a tea.Msg that indicates output text has been received from
// the server.
type OutputReceived struct {
	Output *api.ClientOutput
	Tag    string
}
