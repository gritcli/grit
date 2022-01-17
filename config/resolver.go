package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gritcli/grit/driver/registry"
	"github.com/gritcli/grit/driver/sourcedriver"
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
// normalizes any values that may be specified in multiple ways. It does NOT
// inspect the driver-specific configuration values within "source" blocks.
//
// 2. The "assemble" phase produces a Config struct value from the merged
// configuration.
type resolver struct {
	reg *registry.Registry
	cfg unresolvedConfig
}

// Merge merges the configuration from c.
func (r *resolver) Merge(filename string, c configFile) error {
	if c.DaemonBlock != nil {
		if err := mergeDaemonBlock(&r.cfg, filename, *c.DaemonBlock); err != nil {
			return err
		}
	}

	if c.ClonesDefaultsBlock != nil {
		if err := mergeClonesDefaultsBlock(&r.cfg, filename, *c.ClonesDefaultsBlock); err != nil {
			return err
		}
	}

	if c.GitDefaultsBlock != nil {
		if err := mergeGitDefaultsBlock(&r.cfg, filename, *c.GitDefaultsBlock); err != nil {
			return err
		}
	}

	for _, b := range c.SourceBlocks {
		if err := mergeSourceBlock(
			r.reg,
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
	mergeDefaultSources(r.reg, &r.cfg)

	if err := normalizeDaemonBlock(&r.cfg); err != nil {
		return err
	}

	if err := normalizeClonesDefaultsBlock(&r.cfg); err != nil {
		return err
	}

	if err := normalizeGitDefaultsBlock(&r.cfg); err != nil {
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
func (r *resolver) Assemble() (Config, error) {
	cfg := Config{
		Daemon:         assembleDaemonBlock(r.cfg.Daemon.Block),
		ClonesDefaults: assembleClonesBlock(r.cfg.ClonesDefaults.Block),
		GitDefaults:    assembleGitBlock(r.cfg.GitDefaults.Block),
	}

	for _, s := range r.cfg.Sources {
		src, err := assembleSourceBlock(r.cfg, s)
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

// unresolvedConfig is a configuration that is in the process of being resolved.
type unresolvedConfig struct {
	// Daemon contains information about the first "daemon" block found within
	// the configuration files. Only one of the loaded files may contain a
	// "daemon" block.
	Daemon struct {
		Block daemonBlock
		File  string
	}

	// ClonesDefaults contains information about the first (root-level) "clones"
	// defaults block found within the configuration files. Only one of the
	// loaded files may contain a "clones" defaults block.
	ClonesDefaults struct {
		Block clonesBlock
		File  string
	}

	// GitDefaults contains information about the first (root-level) "git"
	// defaults block found within the configuration files. Only one of the
	// loaded files may contain a "git" defaults block.
	GitDefaults struct {
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
	Block       sourceBlock
	DriverBlock sourcedriver.ConfigSchema
	File        string
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
