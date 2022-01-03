package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	homedir "github.com/mitchellh/go-homedir"
)

// Load loads the configuration from all files in the given directory.
//
// If the directory doesn't exist or does not contain any configuration files,
// then DefaultConfig is returned.
func Load(dir string) (Config, error) {
	var l loader

	if err := l.LoadDir(dir); err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig, nil
		}

		return Config{}, err
	}

	return l.Finalize()
}

// loader loads and assembles a configuration from several configuration files.
type loader struct {
	config             Config
	daemonBlockFile    string
	globalGitBlockFile string
	sourceBlocks       []sourceBlock
	sourceBlockFiles   map[string]string
}

// Finalize returns the loaded configuration.
func (l *loader) Finalize() (Config, error) {
	l.mergeDefaults()

	for _, b := range l.sourceBlocks {
		filename := l.sourceBlockFiles[b.Name]

		if err := l.mergeSourceBlock(filename, b); err != nil {
			return Config{}, fmt.Errorf("%s: %w", filename, err)
		}
	}

	return l.config, nil
}

// LoadDir loads the configuration from all .hcl files in the given directory.
func (l *loader) LoadDir(dir string) error {
	dir, err := homedir.Expand(dir)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Type().IsDir() {
			continue
		}

		name := entry.Name()
		if name[0] == '_' {
			continue
		}

		ext := filepath.Ext(name)
		if ext != ".hcl" {
			continue
		}

		filename := filepath.Join(dir, name)
		if err := l.LoadFile(filename); err != nil {
			return err
		}
	}

	return nil
}

// LoadFile loads the configuration from a single file.
func (l *loader) LoadFile(filename string) error {
	var c configFile
	if err := hclsimple.DecodeFile(filename, nil, &c); err != nil {
		return err
	}

	if c.DaemonBlock != nil {
		if err := l.mergeDaemonBlock(filename, *c.DaemonBlock); err != nil {
			return fmt.Errorf("%s: %w", filename, err)
		}
	}

	if c.GitBlock != nil {
		if err := l.mergeGlobalGitBlock(filename, *c.GitBlock); err != nil {
			return fmt.Errorf("%s: %w", filename, err)
		}
	}

	for _, b := range c.SourceBlocks {
		// Note that source blocks are not merged immediately, instead we
		// perform some cursory validation, and delay merging the final source
		// configuration until after the global defaults have been populated.
		//
		// This is necessary to allow source-specific configuration to fall-back
		// to global configuration options.
		if err := l.prepareSourceBlock(filename, b); err != nil {
			return fmt.Errorf("%s: %w", filename, err)
		}
	}

	return nil
}

// mergeDefaults merges blocks from DefaultConfig that were not explicitly
// defined in the configuration files.
func (l *loader) mergeDefaults() {
	if l.config.Daemon == (Daemon{}) {
		l.config.Daemon = DefaultConfig.Daemon
	}

	if l.config.GlobalGit == (Git{}) {
		l.config.GlobalGit = DefaultConfig.GlobalGit
	}

	if l.config.Sources == nil {
		l.config.Sources = map[string]Source{}
	}

	for n, s := range DefaultConfig.Sources {
		if _, ok := l.config.Sources[n]; !ok {
			l.config.Sources[n] = s
		}
	}
}

// normalizePath expands references to ~ in filesystem names, and resolves
// relative paths.
func normalizePath(filename string, p *string) error {
	n := *p

	n, err := homedir.Expand(n)
	if err != nil {
		return err
	}

	if !filepath.IsAbs(n) {
		dir := filepath.Dir(filename)
		n = filepath.Join(dir, n)
	}

	*p = filepath.Clean(n)

	return nil
}
