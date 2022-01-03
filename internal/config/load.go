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

	return l.Get(), nil
}

// loader loads and assembles a configuration from several configuration files.
type loader struct {
	config      Config
	daemonFile  string
	sourceFiles map[string]string
}

// Get returns the loaded configuration.
func (l *loader) Get() Config {
	l.mergeDefaults()

	return l.config
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
			return err
		}
	}

	for _, sb := range c.SourceBlocks {
		if err := l.mergeSourceBlock(filename, sb); err != nil {
			return err
		}
	}

	return nil
}

// mergeDaemonBlock merges db into the configuration.
func (l *loader) mergeDaemonBlock(filename string, db daemonBlock) error {
	if l.daemonFile != "" {
		return fmt.Errorf("%s: the daemon configuration has already been defined in %s", filename, l.daemonFile)
	}

	d, err := db.resolve(filename)
	if err != nil {
		return fmt.Errorf("%s: the daemon configuration is invalid: %w", filename, err)
	}

	l.daemonFile = filename
	l.config.Daemon = d

	return nil
}

// mergeDaemonBlock merges sb into the configuration.
func (l *loader) mergeSourceBlock(filename string, sb sourceBlock) error {
	if l.config.Sources == nil {
		l.config.Sources = map[string]Source{}
		l.sourceFiles = map[string]string{}
	} else if _, ok := l.config.Sources[sb.Name]; ok {
		return fmt.Errorf(
			"%s: the '%s' repository source has already been defined in %s",
			filename,
			sb.Name,
			l.sourceFiles[sb.Name],
		)
	}

	src, err := sb.resolve(filename)
	if err != nil {
		return fmt.Errorf(
			"%s: the '%s' repository source is invalid: %w",
			filename,
			sb.Name,
			err,
		)
	}

	l.sourceFiles[src.Name] = filename
	l.config.Sources[src.Name] = src

	return nil
}

// mergeDefaults merges blocks from DefaultConfig that were not explicitly
// defined in the configuration files.
func (l *loader) mergeDefaults() {
	if l.daemonFile == "" {
		l.config.Daemon = DefaultConfig.Daemon
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
