package config_test

import (
	. "github.com/gritcli/grit/internal/config"
	homedir "github.com/mitchellh/go-homedir"
)

// defaultConfig is the expected default Grit configuration.
var defaultConfig = Config{
	Daemon: Daemon{
		Socket: "~/grit/daemon.sock",
	},
	ClonesDefaults: Clones{
		Dir: "~/grit",
	},
	Sources: map[string]Source{
		"github": {
			Name:    "github",
			Enabled: true,
			Clones: Clones{
				Dir: "~/grit/github",
			},
			Driver: GitHub{
				Domain: "github.com",
			},
		},
	},
}

func init() {
	var err error
	defaultConfig.Daemon.Socket, err = homedir.Expand(defaultConfig.Daemon.Socket)
	if err != nil {
		panic(err)
	}

	defaultConfig.ClonesDefaults.Dir, err = homedir.Expand(defaultConfig.ClonesDefaults.Dir)
	if err != nil {
		panic(err)
	}

	for n, s := range defaultConfig.Sources {
		s.Clones.Dir, err = homedir.Expand(s.Clones.Dir)
		if err != nil {
			panic(err)
		}

		defaultConfig.Sources[n] = s
	}
}

// withDaemon returns a copy of cfg with a different daemon configuration.
func withDaemon(cfg Config, d Daemon) Config {
	cfg.Daemon = d
	return cfg
}

// withClonesDefaults returns a copy of cfg with a different clones defaults
// configuration.
func withClonesDefaults(cfg Config, c Clones) Config {
	cfg.ClonesDefaults = c
	return cfg
}

// withGitDefaults returns a copy of cfg with a different git defaults
// configuration.
func withGitDefaults(cfg Config, g Git) Config {
	cfg.GitDefaults = g
	return cfg
}

// withSource returns a copy of cfg with an additional repository source.
func withSource(cfg Config, src Source) Config {
	prev := cfg.Sources
	cfg.Sources = map[string]Source{}

	var err error
	src.Clones.Dir, err = homedir.Expand(src.Clones.Dir)
	if err != nil {
		panic(err)
	}

	for n, s := range prev {
		cfg.Sources[n] = s
	}

	cfg.Sources[src.Name] = src

	return cfg
}
