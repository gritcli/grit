package githubsource

import (
	"github.com/gritcli/grit/driver/sourcedriver"
)

// Registration contains information about the driver used to register it with
// Grit's driver registry.
var Registration = sourcedriver.Registration{
	Name:        "github",
	Description: "adds support for GitHub and GitHub Enterprise Server as repository sources",
	NewConfigSchema: func() sourcedriver.ConfigSchema {
		return &configSchema{}
	},
	ImplicitSources: map[string]sourcedriver.ConfigSchema{
		"github": &configSchema{
			Domain: "github.com",
		},
	},
}
