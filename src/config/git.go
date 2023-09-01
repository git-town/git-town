package config

// Git manages configuration data stored in Git metadata.
// Supports configuration in the local repo and the global Git configuration.
type Git struct {
	querierRunner
	config GitConfig
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(runner querierRunner, config GitConfig) Git {
	return Git{
		config:        config,
		querierRunner: runner,
	}
}

func (g Git) GlobalConfigClone() GitConfigCache {
	return g.config.Global.Clone()
}

func (g Git) GlobalConfigValue(key Key) string {
	return g.config.Global[key]
}

func (g Git) LocalConfigClone() GitConfigCache {
	return g.config.Local.Clone()
}

func (g Git) LocalConfigKeysMatching(pattern string) []Key {
	return g.config.Local.KeysMatching(pattern)
}

func (g Git) LocalConfigValue(key Key) string {
	return g.config.Local[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (g Git) LocalOrGlobalConfigValue(key Key) string {
	local := g.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return g.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
// TODO: replace this with assigmnents.
func (g *Git) Reload() {
	g.config.Local = LoadGitConfigCache(g.querierRunner, false)
	g.config.Global = LoadGitConfigCache(g.querierRunner, true)
}

func (g *Git) RemoveGlobalConfigValue(key Key) error {
	delete(g.config.Global, key)
	return g.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (g *Git) RemoveLocalConfigValue(key Key) error {
	delete(g.config.Local, key)
	err := g.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (g *Git) SetGlobalConfigValue(key Key, value string) error {
	g.config.Global[key] = value
	return g.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (g *Git) SetLocalConfigValue(key Key, value string) error {
	g.config.Local[key] = value
	return g.Run("git", "config", key.String(), value)
}
