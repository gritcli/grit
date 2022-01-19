package config

import (
	"os"
	"path/filepath"

	"github.com/gritcli/grit/driver/registry"
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
	dir, err := homedir.Expand(dir)
	if err != nil {
		return Config{}, err
	}

	// TODO: make sure dir is absolute

	r := resolver{
		configDir: dir,
		registry: registry.Registry{
			Parent: reg,
		},
	}

	if err := loadDir(&r, dir); err != nil {
		return Config{}, err
	}

	return r.Finalize()
}

// loadDir loads the configuration from all .hcl files in the given
// directory and merges them into the configuration using r.
func loadDir(r *resolver, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, entry := range entries {
		if entry.Type().IsDir() {
			continue
		}

		name := entry.Name()
		if name[0] == '_' || name[0] == '.' {
			continue
		}

		ext := filepath.Ext(name)
		if ext != ".hcl" {
			continue
		}

		filename := filepath.Join(dir, name)

		var f fileSchema
		if err := hclsimple.DecodeFile(filename, nil, &f); err != nil {
			return err
		}

		if err := r.Merge(filename, f); err != nil {
			return err
		}
	}

	return nil
}
