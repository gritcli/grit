package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gritcli/grit/driver/registry"
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

	for _, b := range c.VCSDefaultsBlocks {
		if err := mergeVCSDefaultsBlock(
			r.reg,
			&r.cfg,
			filename,
			b,
		); err != nil {
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
	if err := mergeImplicitVCSDefaults(r.reg, &r.cfg); err != nil {
		return err
	}

	mergeImplicitSources(r.reg, &r.cfg)

	if err := normalizeDaemonBlock(&r.cfg); err != nil {
		return err
	}

	if err := normalizeClonesDefaultsBlock(&r.cfg); err != nil {
		return err
	}

	for k, s := range r.cfg.Sources {
		if err := normalizeSourceBlock(
			r.reg,
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
	cfg := Config{
		Daemon: assembleDaemonBlock(r.cfg.Daemon.Block),
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

	// VCSDefaults contains information about the "vcs" blocks within the
	// configuration files.
	VCSDefaults map[string]unresolvedVCS

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
