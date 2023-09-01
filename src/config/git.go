package config

// Git manages configuration data stored in Git metadata.
// Supports configuration in the local repo and the global Git configuration.
type Git struct {
	runner

	// globalConfigCache is a cache of the global Git configuration.
	globalConfig GitConfigCache

	// localConfigCache is a cache of the Git configuration in the local Git repo.
	localConfig GitConfigCache
}

type runner interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(runner runner) Git {
	return Git{
		localConfig:  LoadGitConfig(runner, false),
		globalConfig: LoadGitConfig(runner, true),
		runner:       runner,
	}
}

func (g Git) GlobalConfigClone() GitConfigCache {
	return g.globalConfig.Clone()
}

func (g Git) GlobalConfigValue(key Key) string {
	return g.globalConfig[key]
}

func (g Git) LocalConfigClone() GitConfigCache {
	return g.localConfig.Clone()
}

func (g Git) LocalConfigValue(key Key) string {
	return g.localConfig[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (g *Git) LocalOrGlobalConfigValue(key Key) string {
	local, has := g.localConfig[key]
	if has {
		return local
	}
	return g.globalConfig[key]
}

func (g *Git) LocalConfigKeysMatching(pattern string) []Key {
	return g.localConfig.KeysMatching(pattern)
}

// Reload refreshes the cached configuration information.
// TODO: move this somewhere else?
func (g *Git) Reload() {
	g.localConfig = LoadGitConfig(g.runner, false)
	g.globalConfig = LoadGitConfig(g.runner, true)
}

func (g *Git) RemoveGlobalConfigValue(key Key) (string, error) {
	delete(g.globalConfig, key)
	return g.Query("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (g *Git) RemoveLocalConfigValue(key Key) error {
	delete(g.localConfig, key)
	err := g.runner.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (g *Git) SetGlobalConfigValue(key Key, value string) (string, error) {
	g.globalConfig[key] = value
	return g.runner.Query("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (g *Git) SetLocalConfigValue(key Key, value string) error {
	g.localConfig[key] = value
	return g.runner.Run("git", "config", key.String(), value)
}
