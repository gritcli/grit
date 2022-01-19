package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2/hclsimple"
	homedir "github.com/mitchellh/go-homedir"
)

// Load loads the configuration from all files in the given directory.
//
// If the directory doesn't exist or does not contain any configuration files,
// then DefaultConfig is returned.
//
// reg is a registry of drivers that can be used within the configuration. It
// may be nil.
func Load(dir string, reg *registry.Registry) (Config, error) {
	l := loader{
		Registry: registry.Registry{
			Parent: reg,
		},
	}

	if err := l.LoadDir(dir); err != nil {
		return Config{}, err
	}

	return l.Finalize()
}

// A loader loads configuration from files and flattens them into a single
// coherent configuration value.
//
// The loader has two phases:
//
// 1. The "merge" phase merges configuration from separate files into an
// intermediate representation of the configuration.
//
// 2. The "finalize" phase produces a Config value from the intermediate
// representation.
type loader struct {
	// Registry is used to lookup information about source and VCS drivers.
	Registry registry.Registry

	// configDir is the directory containing the configuration files to load.
	//
	// It must be an absolute directory.
	configDir string

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

// LoadDir loads the configuration from .hcl files in the given directory and
// merges them into the configuration.
//
// This is the main entrypoint for the loader.
//
// Files that begin with and underscore (_) or dot (.) are ignored. It does not
// descend into sub-directories.
func (l *loader) LoadDir(dir string) error {
	dir, err := homedir.Expand(dir)
	if err != nil {
		return err
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}

	l.configDir = dir

	entries, err := os.ReadDir(l.configDir)
	if err != nil {
		if os.IsNotExist(err) {
			// A non-existent directory is not an error, it's simply an empty
			// configuration.
			return nil
		}

		return err
	}

	for _, entry := range entries {
		if !isConfigFile(entry) {
			continue
		}

		file := filepath.Join(dir, entry.Name())

		var f fileSchema
		if err := hclsimple.DecodeFile(file, nil, &f); err != nil {
			return err
		}

		if err := l.mergeFile(file, f); err != nil {
			return err
		}
	}

	return nil
}

// isConfigFile returns true if e represents a config file that should be
// loaded.
func isConfigFile(e fs.DirEntry) bool {
	if e.Type().IsDir() {
		return false
	}

	name := e.Name()
	if name[0] == '_' || name[0] == '.' {
		return false
	}

	return strings.EqualFold(
		filepath.Ext(name),
		".hcl",
	)
}

// mergeFile merges the configuration from a single file into the intermediate
// representation of the configuration.
func (l *loader) mergeFile(file string, f fileSchema) error {
	if f.Daemon != nil {
		if err := l.mergeDaemon(file, *f.Daemon); err != nil {
			return err
		}
	}

	if f.GlobalClones != nil {
		if err := l.mergeGlobalClones(file, *f.GlobalClones); err != nil {
			return err
		}
	}

	for _, vcs := range f.GlobalVCSs {
		if err := l.mergeGlobalVCS(file, vcs); err != nil {
			return err
		}
	}

	for _, s := range f.Sources {
		if err := l.mergeSource(file, s); err != nil {
			return err
		}
	}

	return nil
}

// Finalize builds the complete configuration from the intermediate
// representation.
func (l *loader) Finalize() (Config, error) {
	if err := l.populateDaemonDefaults(); err != nil {
		return Config{}, err
	}

	if err := l.populateGlobalClonesDefaults(); err != nil {
		return Config{}, err
	}

	if err := l.populateImplicitGlobalVCSs(); err != nil {
		return Config{}, err
	}

	l.populateImplicitSources()

	cfg := Config{
		Daemon: l.daemon,
	}

	for _, i := range l.sources {
		src, err := l.finalizeSource(i)
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

// normalizePath sets *p to its absolute path representation, relative to the
// configuration directory.
//
// Leading tildes (~) are expanded to the user's home directory.
//
// It does not require the referenced path to exist, and hence does not follow
// symlinks, etc.
//
// It does nothing if p is nil or *p is empty.
func (l *loader) normalizePath(p *string) error {
	if p == nil || *p == "" {
		return nil
	}

	result, err := homedir.Expand(*p)
	if err != nil {
		return err
	}

	if !filepath.IsAbs(result) {
		result = filepath.Join(l.configDir, result)
	}

	*p = filepath.Clean(result)

	return nil
}
