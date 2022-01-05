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
	GlobalGit: Git{
		PrivateKey: "",
		PreferHTTP: false,
	},
	Sources: map[string]Source{
		"github": {
			Name:    "github",
			Enabled: true,
			Config: GitHubConfig{
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
}

// withDaemon returns a copy of cfg with a different daemon configuration.
func withDaemon(cfg Config, d Daemon) Config {
	cfg.Daemon = d
	return cfg
}

// withGlobalGit returns a copy of cfg with a different git configuration.
func withGlobalGit(cfg Config, g Git) Config {
	cfg.GlobalGit = g
	return cfg
}

// withSource returns a copy of cfg with an additional repository source.
func withSource(cfg Config, src Source) Config {
	prev := cfg.Sources
	cfg.Sources = map[string]Source{}

	for n, s := range prev {
		cfg.Sources[n] = s
	}

	cfg.Sources[src.Name] = src

	return cfg
}
