package githubsource

import "github.com/gritcli/grit/daemon/internal/driver/sourcedriver"

// Registration contains information about the driver used to register it with
// Grit's driver registry.
var Registration = sourcedriver.Registration{
	Name:         "github",
	Description:  "adds support for GitHub and GitHub Enterprise Server as repository sources",
	ConfigLoader: configLoader{},
}
