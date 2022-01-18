package config

import "github.com/hashicorp/hcl/v2"

// fileSchema is HCL schema for a single configuration file.
type fileSchema struct {
	// Daemon is the configuratio for the Grit daemon.
	//
	// It is optional, but if present can only be specified in one file.
	Daemon *daemonSchema `hcl:"daemon,block"`

	// GlobalClones is the global (non-source-specific) configuration for how
	// Grit stores local repository clones.
	//
	// It is optional, but if present can only be specified in one file.
	GlobalClones *clonesSchema `hcl:"clones,block"`

	VCSDefaults []vcsSchema    `hcl:"vcs,block"`
	Sources     []sourceSchema `hcl:"source,block"`
}

// daemonSchema is the HCL schema for a "daemon" block.
type daemonSchema struct {
	// Socket is the path to the unix-socket address used for gRPC communication
	// between the CLI and the daemon.
	Socket string `hcl:"socket,optional"`
}

// clonesSchema is the HCL schema for a "clones" block.
type clonesSchema struct {
	// Dir is the directory in which local repository clones are stored.
	//
	// If this is a global clones block (non-source-specific), this is the base
	// directory under which each source has its own directory by default.
	//
	// For clones configuration within a specific source, this is the exact
	// path under which clones are stored.
	Dir string `hcl:"dir,optional"`
}

// vcsSchema is the HCL schema for a "vcs" block.
type vcsSchema struct {
	DriverAlias string   `hcl:",label"`
	Body        hcl.Body `hcl:",remain"`
}

// sourceSchema is the HCL schema for a "source" block.
type sourceSchema struct {
	Name        string        `hcl:",label"`
	DriverAlias string        `hcl:",label"`
	Enabled     *bool         `hcl:"enabled"`
	ClonesBlock *clonesSchema `hcl:"clones,block"`
	VCSBlocks   []vcsSchema   `hcl:"vcs,block"`
	DriverBlock hcl.Body      `hcl:",remain"`
}
