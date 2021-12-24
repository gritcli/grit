package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Load loads the configuration from all files in the given directory.
func Load(dir string) (Config, error) {
	var cfg Config

	entries, err := os.ReadDir(dir)
	if err != nil {
		return Config{}, err
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

		if name == "grit.hcl" {
			err = loadMainConfig(filename, &cfg)
		} else if strings.HasSuffix(name, ".source.hcl") {
			err = loadSourceConfig(filename, &cfg)
		} else {
			err = fmt.Errorf("%s: unrecognized configuration file", filename)
		}

		if err != nil {
			return Config{}, err
		}
	}

	if err := mergeDefaults(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
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
