package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	homedir "github.com/mitchellh/go-homedir"
)

// Load loads the configuration from all files in the given directory.
//
// If the directory doesn't exist DefaultConfig is returned.
func Load(dir string) (Config, error) {
	var cfg Config

	if err := scanDir(&cfg, dir); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// scanDir scans the contents of dir and loads any configuration files found
// within.
func scanDir(cfg *Config, dir string) error {
	dir, err := homedir.Expand(dir)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			*cfg = DefaultConfig
			return nil
		}

		return err
	}

	for _, entry := range entries {
		if err := loadEntry(cfg, dir, entry); err != nil {
			return err
		}
	}

	if err := mergeDefaults(cfg); err != nil {
		return err
	}

	if err := normalize(cfg, dir); err != nil {
		return err
	}

	return nil
}

// loadEntry loads configuration from a directory entry if it refers to a known
// .hcl file.
func loadEntry(cfg *Config, dir string, entry os.DirEntry) error {
	if entry.Type().IsDir() {
		return nil
	}

	name := entry.Name()
	if name[0] == '_' {
		return nil
	}

	ext := filepath.Ext(name)
	if ext != ".hcl" {
		return nil
	}

	filename := filepath.Join(dir, name)

	if name == "grit.hcl" {
		return loadMainConfig(filename, cfg)
	}

	if strings.HasSuffix(name, ".source.hcl") {
		return loadSourceConfig(filename, cfg)
	}

	return fmt.Errorf("%s: unrecognized configuration file", filename)
}

// loadMainConfig loads the config in the "grit.hcl" file.
func loadMainConfig(filename string, cfg *Config) error {
	var c mainConfig
	if err := hclsimple.DecodeFile(filename, nil, &c); err != nil {
		return err
	}

	if c.Daemon != nil {
		cfg.Daemon = *c.Daemon
	}

	return nil
}

// loadMainConfig loads the config in a "*.source.hcl" file.
func loadSourceConfig(filename string, cfg *Config) error {
	var c sourceConfig
	if err := hclsimple.DecodeFile(filename, nil, &c); err != nil {
		return err
	}

	return nil
}

// mergeDefaults merges default values into cfg.
func mergeDefaults(cfg *Config) error {
	if cfg.Daemon.Socket == "" {
		cfg.Daemon.Socket = DefaultConfig.Daemon.Socket
	}

	return nil
}

// normalize normalizes the structure and values within cfg.
func normalize(cfg *Config, configDir string) error {
	return normalizePath(&cfg.Daemon.Socket, configDir)
}

// expandPath expands references to ~ in filesystem names.
func normalizePath(p *string, configDir string) error {
	n, err := homedir.Expand(*p)
	if err != nil {
		return err
	}

	if !filepath.IsAbs(n) {
		n = filepath.Join(configDir, n)
	}

	*p = filepath.Clean(n)

	return nil
}
