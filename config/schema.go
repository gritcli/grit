package config

import "github.com/hashicorp/hcl/v2"

// fileSchema is HCL schema for a single configuration file.
type fileSchema struct {
	Daemon         *daemonSchema  `hcl:"daemon,block"`
	ClonesDefaults *clonesSchema  `hcl:"clones,block"`
	VCSDefaults    []vcsSchema    `hcl:"vcs,block"`
	Sources        []sourceSchema `hcl:"source,block"`
}

// daemonSchema is the HCL schema for a "daemon" block.
type daemonSchema struct {
	Socket string `hcl:"socket,optional"`
}

// clonesSchema is the HCL schema for a "clones" block.
type clonesSchema struct {
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
