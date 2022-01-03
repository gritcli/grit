package config

import "fmt"

// mergeGlobalGitBlock merges b into the configuration.
func (l *loader) mergeGlobalGitBlock(filename string, b gitBlock) error {
	if l.globalGitBlockFile != "" {
		return fmt.Errorf("the global git configuration has already been defined in %s", l.globalGitBlockFile)
	}

	g := Git(b)

	if g.PrivateKey != "" {
		if err := normalizePath(filename, &g.PrivateKey); err != nil {
			return fmt.Errorf("the git private key path can not be resolved: %w", err)
		}
	}

	l.globalGitBlockFile = filename
	l.config.GlobalGit = g

	return nil
}
