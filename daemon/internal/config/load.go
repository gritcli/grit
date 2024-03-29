package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gritcli/grit/daemon/internal/driver/vcsdriver"
	"github.com/hashicorp/hcl/v2"
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
func Load(dir string, reg *DriverRegistry) (Config, error) {
	if dir == "" {
		dir = DefaultDirectory
	}

	configDir, err := homedir.Expand(dir)
	if err != nil {
		return Config{}, fmt.Errorf(
			"unable to resolve configuration directory: %w (%s)",
			err,
			dir,
		)
	}

	configDir, err = filepath.Abs(configDir)
	if err != nil {
		// CODE COVERAGE: I'm not aware of a simple and cross-platform way to
		// induce filepath.Abs() to fail.
		return Config{}, err
	}

	l := loader{
		ConfigDir: configDir,
		Registry: DriverRegistry{
			Parent: reg,
		},
	}

	return l.Load()
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
	// ConfigDir is the directory containing the configuration files to load.
	//
	// It must be an absolute directory.
	ConfigDir string

	// Registry is used to lookup information about source and VCS drivers.
	Registry DriverRegistry

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

	// defaultVCSs is a map of VCS driver name to the default configuration for
	// that driver.
	defaultVCSs map[string]vcsdriver.Config

	// globalVCSs is a map of VCS driver name to the global configuration for
	// that driver.
	globalVCSs map[string]vcsdriver.Config

	// sources is a map of _lowercase_ source name to its intermediate
	// representation.
	sources map[string]intermediateSource
}

// Load loads the configuration from .hcl files in l.ConfigDir and merges them
// into the configuration.
//
// This is the main entrypoint for the loader.
//
// Files that begin with and underscore (_) or dot (.) are ignored. It does not
// descend into sub-directories.
func (l *loader) Load() (Config, error) {
	if err := l.prepare(); err != nil {
		return Config{}, err
	}

	if err := l.load(); err != nil {
		return Config{}, err
	}

	return l.finalize()
}

// prepare prepares the loader to load configuration.
func (l *loader) prepare() error {
	if err := l.populateDefaultVCSs(); err != nil {
		return err
	}

	return nil
}

func (l *loader) load() error {
	entries, err := os.ReadDir(l.ConfigDir)
	if err != nil {
		if os.IsNotExist(err) {
			// A non-existent directory is not an error, it's simply an empty
			// configuration.
			err = nil
		}

		return err
	}

	for _, entry := range entries {
		if !isConfigFile(entry) {
			continue
		}

		file := filepath.Join(l.ConfigDir, entry.Name())

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
func (l *loader) mergeFile(file string, f fileSchema) (err error) {
	defer func() {
		if err != nil && !isHCLError(err) {
			err = fmt.Errorf("%s: %w", file, err)
		}
	}()

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

// finalize builds the complete configuration from the intermediate
// representation.
func (l *loader) finalize() (Config, error) {
	if err := l.populateDaemonDefaults(); err != nil {
		return Config{}, err
	}

	if err := l.populateGlobalClonesDefaults(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Daemon: l.daemon,
	}

	if err := l.populateImplicitSources(&cfg); err != nil {
		return Config{}, err
	}

	for _, i := range l.sources {
		src, err := l.finalizeSource(i)
		if err != nil {
			if i.File != "" && !isHCLError(err) {
				return Config{}, fmt.Errorf("%s: %w", i.File, err)
			}

			return Config{}, err
		}

		cfg.Sources = append(cfg.Sources, src)
	}

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
		result = filepath.Join(l.ConfigDir, result)
	}

	*p = filepath.Clean(result)

	return nil
}

// isHCLError returns true if err is an hcl.Diagnostics error, in which case it
// will already contain file/line information.
func isHCLError(err error) bool {
	var diag hcl.Diagnostics
	return errors.As(err, &diag)
}
