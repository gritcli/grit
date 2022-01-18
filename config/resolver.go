package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/vcsdriver"
	homedir "github.com/mitchellh/go-homedir"
)

// A resolver flattens multiple configuration files into a single coherent
// configuration.
type resolver struct {
	// configDir is the directory containing the configuration files to load.
	configDir string

	registry *registry.Registry
	cfg      unresolvedConfig

	// output is the configuration that is built by the resolver.
	output Config

	// currentFile is the name of the file currently being merged.
	currentFile string

	// daemonFile is the name of the file containing the "daemon" block that was
	// merged into the configuration. It is empty if no "daemon" block has yet
	// been merged.
	daemonFile string

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
}

// Merge merges the configuration from c.
func (r *resolver) Merge(filename string, f fileSchema) error {
	r.currentFile = filename
	defer func() {
		r.currentFile = ""
	}()

	if f.Daemon != nil {
		if err := r.mergeDaemon(*f.Daemon); err != nil {
			return err
		}
	}

	if f.GlobalClones != nil {
		if err := r.mergeGlobalClones(*f.GlobalClones); err != nil {
			return err
		}
	}

	for _, vcs := range f.GlobalVCSs {
		if err := r.mergeGlobalVCS(vcs); err != nil {
			return err
		}
	}

	for _, b := range f.Sources {
		if err := mergeSourceBlock(
			r.registry,
			&r.cfg,
			filename,
			b,
		); err != nil {
			return err
		}
	}

	return nil
}

// Normalize normalizes the configuration and populates it with default values.
func (r *resolver) Normalize() error {
	if err := r.populateDaemonDefaults(); err != nil {
		return err
	}

	if err := r.populateGlobalClonesDefaults(); err != nil {
		return err
	}

	if err := r.populateImplicitGlobalVCSs(); err != nil {
		return err
	}

	mergeImplicitSources(r.registry, &r.cfg)

	for k, s := range r.cfg.Sources {
		if err := normalizeSourceBlock(
			r,
			r.registry,
			r.cfg,
			&s,
		); err != nil {
			return err
		}

		r.cfg.Sources[k] = s
	}

	return nil
}

// Assemble returns the file configuration assembled from the various source
// files.
func (r *resolver) Assemble() (Config, error) {
	for _, s := range r.cfg.Sources {
		src, err := assembleSourceBlock(r, r.cfg, s)
		if err != nil {
			return Config{}, err
		}

		r.output.Sources = append(r.output.Sources, src)
	}

	sort.Slice(
		r.output.Sources,
		func(i, j int) bool {
			return r.output.Sources[i].Name < r.output.Sources[j].Name
		},
	)

	return r.output, nil
}

// unresolvedConfig is a configuration that is in the process of being resolved.
type unresolvedConfig struct {
	// Sources contains information about the "source" blocks within the
	// configuration files.
	Sources map[string]unresolvedSource
}

// normalizePath normalizes the path *p relative to the config file that
// contains it.
//
// If *p begins with a tilde (~), it is resolved relative to the user's home
// directory.
//
// If *p is a relative path, it is resolved to an absolute path relative to the
// directory of the given filename.
//
// It does nothing if p is nil or *p is empty.
//
// It panics if *p is relative and filename is empty.
func normalizePath(filename string, p *string) error {
	if p == nil || *p == "" {
		return nil
	}

	n := *p

	n, err := homedir.Expand(n)
	if err != nil {
		err = fmt.Errorf(
			"unable to expand %s with the user's home directory: %w",
			n,
			err,
		)

		if filename != "" {
			err = fmt.Errorf(
				"%s: %w",
				filename,
				err,
			)
		}

		return err
	}

	if !filepath.IsAbs(n) {
		if filename == "" {
			panic("cannot resolve relative path outside of configuration file")
		}

		dir := filepath.Dir(filename)

		if !filepath.IsAbs(dir) {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			dir = filepath.Join(cwd, dir)
		}

		n = filepath.Join(dir, n)
	}

	*p = filepath.Clean(n)

	return nil
}

// normalizeContext is a (partial) implementation of the ConfigNormalizeContext
// interface from the sourcedriver and vcsdriver packages.
type normalizeContext struct {
	// configDir is the directory containing the config files being loaded.
	configDir string
}

// NormalizePath resolves *p to an absolute path relative to the configuration
// directory. It supports expanding leading tilde (~) to the user's home
// directory.
//
// It does not require the referenced path to exist, and hence does not resolve
// symlinks, etc.
//
// It does nothing if p is nil or *p is empty.
func (n *normalizeContext) NormalizePath(p *string) error {
	if p == nil || *p == "" {
		return nil
	}

	result, err := homedir.Expand(*p)
	if err != nil {
		return err
	}

	if !filepath.IsAbs(result) {
		baseDir, err := n.baseDir()
		if err != nil {
			return err
		}

		result = filepath.Join(baseDir, result)
	}

	*p = filepath.Clean(result)

	return nil
}

// baseDir returns the directory against which relative paths are resolved.
func (n *normalizeContext) baseDir() (string, error) {
	if filepath.IsAbs(n.configDir) {
		return n.configDir, nil
	}

	cwd, err := os.Getwd()
	return filepath.Join(cwd, n.configDir), err
}
