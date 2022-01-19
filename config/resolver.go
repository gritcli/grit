package config

import (
	"path/filepath"
	"sort"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/vcsdriver"
	homedir "github.com/mitchellh/go-homedir"
)

// A resolver flattens multiple configuration files into a single coherent
// configuration.
//
// There are two phases:
//
// 1. The "merge" phase merges configuration from separate files into an
// intermediate representation of the configuration.
//
// 2. The "finalize" phase produces a Config value from the intermediate
// representation.
type resolver struct {
	// configDir is the directory containing the configuration files to load.
	//
	// It must be an absolute directory.
	configDir string

	// registry is used to lookup information about source and VCS drivers.
	registry registry.Registry

	// daemonFile is the name of the file containing the "daemon" block that was
	// merged into the configuration. It is empty if no "daemon" block has yet
	// been merged.
	daemonFile string

	// daemon is the daemon configuration parsed from daemonFile.
	daemon Daemon

	// globalClonesFile is the name of the file containing the global clones
	// configuration. It is empty if no global clones block has been parsed yet.
	globalClonesFile string

	// globalClones is the clones configuration parsed from globalClonesFile.
	globalClones Clones

	// globalVCSFiles is a map of VCS driver name to the file containing the
	// global configuration for that driver.
	globalVCSFiles map[string]string

	// globalVCSs is a map of VCS driver name to the global configuration for
	// that driver.
	globalVCSs map[string]vcsdriver.Config

	// sources is a map of _lowercase_ source name to its intermediate
	// representation.
	sources map[string]intermediateSource
}

// Merge merges the configuration from a single file into the intermediate
// representation of the configuration.
func (r *resolver) Merge(file string, f fileSchema) error {
	if f.Daemon != nil {
		if err := r.mergeDaemon(file, *f.Daemon); err != nil {
			return err
		}
	}

	if f.GlobalClones != nil {
		if err := r.mergeGlobalClones(file, *f.GlobalClones); err != nil {
			return err
		}
	}

	for _, vcs := range f.GlobalVCSs {
		if err := r.mergeGlobalVCS(file, vcs); err != nil {
			return err
		}
	}

	for _, s := range f.Sources {
		if err := r.mergeSource(file, s); err != nil {
			return err
		}
	}

	return nil
}

// Finalize builds the complete configuration from the intermediate
// representation.
func (r *resolver) Finalize() (Config, error) {
	if err := r.populateDaemonDefaults(); err != nil {
		return Config{}, err
	}

	if err := r.populateGlobalClonesDefaults(); err != nil {
		return Config{}, err
	}

	if err := r.populateImplicitGlobalVCSs(); err != nil {
		return Config{}, err
	}

	r.populateImplicitSources()

	cfg := Config{
		Daemon: r.daemon,
	}

	for _, i := range r.sources {
		src, err := r.finalizeSource(i)
		if err != nil {
			return Config{}, err
		}

		cfg.Sources = append(cfg.Sources, src)
	}

	sort.Slice(
		cfg.Sources,
		func(i, j int) bool {
			return cfg.Sources[i].Name < cfg.Sources[j].Name
		},
	)

	return cfg, nil
}

// normalizePath resolves *p to an absolute path relative to the configuration
// directory. It supports expanding leading tilde (~) to the user's home
// directory.
//
// It does not require the referenced path to exist, and hence does not resolve
// symlinks, etc.
//
// It does nothing if p is nil or *p is empty.
func (r *resolver) normalizePath(p *string) error {
	if p == nil || *p == "" {
		return nil
	}

	result, err := homedir.Expand(*p)
	if err != nil {
		return err
	}

	if !filepath.IsAbs(result) {
		result = filepath.Join(r.configDir, result)
	}

	*p = filepath.Clean(result)

	return nil
}
