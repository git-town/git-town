package config

// Git manages configuration data stored in Git metadata.
// Supports configuration in the local repo and the global Git configuration.
type Git struct {
	runner
	config GitConfig
}

type runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(gitConfig GitConfig, runner runner) Git {
	return Git{
		config: gitConfig,
		runner: runner,
	}
}

func (self Git) GlobalConfigClone() GitConfigCache {
	return self.config.Global.Clone()
}

func (self Git) GlobalConfigValue(key Key) string {
	return self.config.Global[key]
}

func (self Git) LocalConfigClone() GitConfigCache {
	return self.config.Local.Clone()
}

func (self Git) LocalConfigKeysMatching(pattern string) []Key {
	return self.config.Local.KeysMatching(pattern)
}

func (self Git) LocalConfigValue(key Key) string {
	return self.config.Local[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (self Git) LocalOrGlobalConfigValue(key Key) string {
	local := self.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return self.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (self *Git) Reload() {
	self.config = LoadGitConfig(self.runner)
}

func (self *Git) RemoveGlobalConfigValue(key Key) error {
	delete(self.config.Global, key)
	return self.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Git) RemoveLocalConfigValue(key Key) error {
	delete(self.config.Local, key)
	err := self.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Git) SetGlobalConfigValue(key Key, value string) error {
	self.config.Global[key] = value
	return self.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Git) SetLocalConfigValue(key Key, value string) error {
	self.config.Local[key] = value
	return self.Run("git", "config", key.String(), value)
}
