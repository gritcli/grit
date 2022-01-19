package githubsource

import (
	"github.com/gritcli/grit/driver/sourcedriver"
)

// Registration contains information about the driver used to register it with
// Grit's driver registry.
var Registration = sourcedriver.Registration{
	Name:        "github",
	Description: "Use repositories from GitHub.com or GitHub Enterprise Server.",
	NewConfigSchema: func() sourcedriver.ConfigSchema {
		return &configSchema{}
	},
	ImplicitSources: map[string]func() sourcedriver.ConfigSchema{
		"github": func() sourcedriver.ConfigSchema {
			return &configSchema{
				Domain: "github.com",
			}
		},
	},
}
