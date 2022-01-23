package gitvcs

import "github.com/gritcli/grit/driver/vcsdriver"

// Registration contains information about the driver used to register it with
// Grit's driver registry.
var Registration = vcsdriver.Registration{
	Name:        "git",
	Description: "adds support for Git repositories",
	NewConfigSchema: func() vcsdriver.ConfigSchema {
		return &configSchema{}
	},
}