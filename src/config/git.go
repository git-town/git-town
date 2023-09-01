package config

// Git manages configuration data stored in Git metadata.
// Supports configuration in the local repo and the global Git configuration.
type Git struct {
	runner

	// globalConfigCache is a cache of the global Git configuration.
	globalConfigCache GitConfigCache

	// localConfigCache is a cache of the Git configuration in the local Git repo.
	localConfigCache GitConfigCache
}

type runner interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(runner runner) Git {
	return Git{
		localConfigCache:  LoadGitConfig(runner, false),
		globalConfigCache: LoadGitConfig(runner, true),
		runner:            runner,
	}
}

func (g Git) GlobalConfigClone() GitConfigCache {
	return g.globalConfigCache.Clone()
}

func (g Git) GlobalConfigValue(key Key) string {
	return g.globalConfigCache[key]
}

func (g Git) LocalConfigClone() GitConfigCache {
	return g.localConfigCache.Clone()
}

func (g Git) LocalConfigValue(key Key) string {
	return g.localConfigCache[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (g *Git) LocalOrGlobalConfigValue(key Key) string {
	local := g.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return g.globalConfigCache[key]
}

func (g *Git) LocalConfigKeysMatching(pattern string) []Key {
	return g.localConfigCache.KeysMatching(pattern)
}

// Reload refreshes the cached configuration information.
func (g *Git) Reload() {
	g.localConfigCache = LoadGitConfig(g.runner, false)
	g.globalConfigCache = LoadGitConfig(g.runner, true)
}

func (g *Git) RemoveGlobalConfigValue(key Key) (string, error) {
	delete(g.globalConfigCache, key)
	return g.Query("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (g *Git) RemoveLocalConfigValue(key Key) error {
	delete(g.localConfigCache, key)
	err := g.runner.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (g *Git) SetGlobalConfigValue(key Key, value string) (string, error) {
	g.globalConfigCache[key] = value
	return g.runner.Query("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (g *Git) SetLocalConfigValue(key Key, value string) error {
	g.localConfigCache[key] = value
	return g.runner.Run("git", "config", key.String(), value)
}
