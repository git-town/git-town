package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// Git manages configuration data stored in Git metadata.
// Supports configuration in the local repo and the global Git configuration.
type Git struct {
	Runner
	GitConfig
}

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(gitConfig GitConfig, runner Runner) Git {
	return Git{
		GitConfig: gitConfig,
		Runner:    runner,
	}
}

func (self Git) GlobalConfigClone() Cache {
	return self.GitConfig.GlobalCache.Clone()
}

func (self Git) GlobalConfigValue(key configdomain.Key) string {
	return self.GitConfig.GlobalCache[key]
}

func (self Git) LocalConfigClone() Cache {
	return self.GitConfig.LocalCache.Clone()
}

func (self Git) LocalConfigKeysMatching(pattern string) []configdomain.Key {
	return self.GitConfig.LocalCache.KeysMatching(pattern)
}

func (self Git) LocalConfigValue(key configdomain.Key) string {
	return self.GitConfig.LocalCache[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (self Git) LocalOrGlobalConfigValue(key configdomain.Key) string {
	local := self.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return self.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (self *Git) Reload() {
	self.GitConfig = LoadGitConfig(self.Runner)
}

func (self *Git) RemoveGlobalConfigValue(key configdomain.Key) error {
	delete(self.GitConfig.GlobalCache, key)
	return self.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Git) RemoveLocalConfigValue(key configdomain.Key) error {
	delete(self.GitConfig.LocalCache, key)
	err := self.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Git) SetGlobalConfigValue(key configdomain.Key, value string) error {
	self.GitConfig.GlobalCache[key] = value
	return self.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Git) SetLocalConfigValue(key configdomain.Key, value string) error {
	self.GitConfig.LocalCache[key] = value
	return self.Run("git", "config", key.String(), value)
}
