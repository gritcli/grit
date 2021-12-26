package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	homedir "github.com/mitchellh/go-homedir"
)

// Load loads the configuration from all files in the given directory.
//
// If the directory doesn't exist DefaultConfig is returned.
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
	l.applyDefaults()
	return l.config
}

// LoadFile loads the configuration from a single file.
func (l *loader) LoadFile(filename string) error {
	var c configFile
	if err := hclsimple.DecodeFile(filename, nil, &c); err != nil {
		return err
	}

	if err := normalize(filename, &c); err != nil {
		return err
	}

	if c.Daemon != nil {
		if l.daemonFile != "" {
			return fmt.Errorf("%s: the daemon configuration has already been defined in %s", filename, l.daemonFile)
		}

		l.daemonFile = filename
		l.config.Daemon = *c.Daemon

	}

	for _, s := range c.Sources {
		if err := l.loadSource(filename, s); err != nil {
			return err
		}
	}

	return nil
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

// loadSource loads a source that was parsed in a configuration file.
func (l *loader) loadSource(filename string, s anySource) error {
	cfg, err := l.decodeDriverConfig(filename, s)
	if err != nil {
		return err
	}

	if l.config.Sources == nil {
		l.config.Sources = map[string]Source{}
		l.sourceFiles = map[string]string{}
	} else if _, ok := l.config.Sources[s.Name]; ok {
		return fmt.Errorf("%s: the '%s' repository source has already been defined in %s", filename, s.Name, l.sourceFiles[s.Name])
	}

	l.sourceFiles[s.Name] = filename
	l.config.Sources[s.Name] = Source{
		Name:   s.Name,
		Config: cfg,
	}

	return nil
}

// decodeDriverConfig decodes a source's configuration using the appropriate
// driver-specific configuration structure.
func (l *loader) decodeDriverConfig(filename string, s anySource) (DriverConfig, error) {
	p, ok := driverConfigPrototypes[s.Driver]
	if !ok {
		return nil, fmt.Errorf("%s: unrecognized source driver: %s", filename, s.Driver)
	}

	ptr := reflect.New(reflect.TypeOf(p))

	diag := gohcl.DecodeBody(s.Body, nil, ptr.Interface())
	if diag.HasErrors() {
		return nil, fmt.Errorf("%s: %w", filename, diag)
	}

	return ptr.Elem().Interface().(DriverConfig), nil
}

// applyDefaults merges missing values from DefaultConfig into cfg.
func (l *loader) applyDefaults() {
	if l.config.Daemon.Socket == "" {
		l.config.Daemon.Socket = DefaultConfig.Daemon.Socket
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

// configFile is the structure of a configuration file as parsed by the HCL
// library.
type configFile struct {
	Daemon  *Daemon     `hcl:"daemon,block"`
	Sources []anySource `hcl:"source,block"`
}

// anySource is a source block that has not been fully parsed.
type anySource struct {
	Name   string   `hcl:",label"`
	Driver string   `hcl:",label"`
	Body   hcl.Body `hcl:",remain"`
}

func normalize(filename string, cfg *configFile) error {
	if cfg.Daemon != nil {
		return normalizePath(filename, &cfg.Daemon.Socket)
	}

	return nil
}

// normalizePath expands references to ~ in filesystem names.
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
