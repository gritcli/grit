package config

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// A resolver flattens multiple configuration files into a single coherent
// configuration.
//
// Resolution is performed in three distinct phases:
//
// 1. The "merge" phase enumerates the content of the parsed files to produce a
// single consistent view of the configuration as though it were contained in a
// single file. This phase is responsible for basic validation and detecting
// semantic errors such as duplicate definitions.
//
// 2. The "normalize" phase sets any optional values to their defaults, and
// normalizes any values that may be specified in multiple ways. This phase
// allows parts of the configuration to depend on other parts of the
// configuration; care must be taken not to introduce circular dependencies
// within the population logic.
//
// 2. The "assemble" phase produces a Config struct value from the merged
// configuration.
type resolver struct {
	cfg unresolvedConfig
}

// Merge merges the configuration from c.
func (r *resolver) Merge(filename string, c configFile) error {
	if c.DaemonBlock != nil {
		if err := mergeDaemonBlock(&r.cfg, filename, *c.DaemonBlock); err != nil {
			return err
		}
	}

	if c.GitBlock != nil {
		if err := mergeGlobalGitBlock(&r.cfg, filename, *c.GitBlock); err != nil {
			return err
		}
	}

	for _, b := range c.SourceBlocks {
		if err := mergeSourceBlock(&r.cfg, filename, b); err != nil {
			return err
		}
	}

	return nil
}

// Normalize normalizes the configuration and populates it with default values.
func (r *resolver) Normalize() error {
	mergeDefaultSources(&r.cfg)

	if err := normalizeDaemonBlock(&r.cfg); err != nil {
		return err
	}

	if err := normalizeGlobalGitBlock(&r.cfg); err != nil {
		return err
	}

	for k, s := range r.cfg.Sources {
		if err := normalizeSourceBlock(r.cfg, &s); err != nil {
			return err
		}

		r.cfg.Sources[k] = s
	}

	return nil
}

// Assemble returns the file configuration assembled from the various source
// files.
func (r *resolver) Assemble() Config {
	cfg := Config{
		Daemon:    assembleDaemonBlock(r.cfg.Daemon.Block),
		GlobalGit: assembleGitBlock(r.cfg.GlobalGit.Block),
		Sources:   map[string]Source{},
	}

	for _, s := range r.cfg.Sources {
		cfg.Sources[s.Block.Name] = assembleSourceBlock(s.Block, s.Body)
	}

	return cfg
}

// unresolvedConfig is a configuration that is in the process of being resolved.
type unresolvedConfig struct {
	// Daemon contains information about the first "daemon" block found within
	// the configuration files. Only one of the loaded files may contain a
	// "daemon" block.
	Daemon struct {
		Block daemonBlock
		File  string
	}

	// GlobalGit contains information about the first global (root-level) "git"
	// block found within the configuration files. Only one of the loaded files
	// may contain a global "git" block.
	GlobalGit struct {
		Block gitBlock
		File  string
	}

	// Sources contains information about a "source" block within the
	// configuration files.
	Sources map[string]unresolvedSource
}

// unresolvedSource contains information about a "source" block within an
// as-yet-unresolved configuration.
type unresolvedSource struct {
	Block sourceBlock
	Body  sourceBlockBody
	File  string
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
func normalizePath(filename string, p *string) error {
	if p == nil || *p == "" {
		return nil
	}

	n := *p

	n, err := homedir.Expand(n)
	if err != nil {
		return fmt.Errorf(
			"%s: unable to expand %s with the user's home directory: %w",
			filename,
			n,
			err,
		)
	}

	if !filepath.IsAbs(n) {
		dir := filepath.Dir(filename)
		n = filepath.Join(dir, n)
	}

	*p = filepath.Clean(n)

	return nil
}