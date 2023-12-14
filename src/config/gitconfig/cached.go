package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// Cached manages configuration data stored in Cached metadata.
// Supports configuration in the local repo and the global Cached configuration.
type Cached struct {
	Runner
	GitConfig
}

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(gitConfig GitConfig, runner Runner) Cached {
	return Cached{
		GitConfig: gitConfig,
		Runner:    runner,
	}
}

func (self Cached) GlobalConfigClone() Cache {
	return self.GitConfig.Global.Clone()
}

func (self Cached) GlobalConfigValue(key configdomain.Key) string {
	return self.GitConfig.Global[key]
}

func (self Cached) LocalConfigClone() Cache {
	return self.GitConfig.Local.Clone()
}

func (self Cached) LocalConfigKeysMatching(pattern string) []configdomain.Key {
	return self.GitConfig.Local.KeysMatching(pattern)
}

func (self Cached) LocalConfigValue(key configdomain.Key) string {
	return self.GitConfig.Local[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (self Cached) LocalOrGlobalConfigValue(key configdomain.Key) string {
	local := self.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return self.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (self *Cached) Reload() {
	self.GitConfig = LoadGitConfig(self.Runner)
}

func (self *Cached) RemoveGlobalConfigValue(key configdomain.Key) error {
	delete(self.GitConfig.Global, key)
	return self.Run("git", "config", "--global", "--unset", key.String())
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Cached) RemoveLocalConfigValue(key configdomain.Key) error {
	delete(self.GitConfig.Local, key)
	err := self.Run("git", "config", "--unset", key.String())
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Cached) SetGlobalConfigValue(key configdomain.Key, value string) error {
	self.GitConfig.Global[key] = value
	return self.Run("git", "config", "--global", key.String(), value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Cached) SetLocalConfigValue(key configdomain.Key, value string) error {
	self.GitConfig.Local[key] = value
	return self.Run("git", "config", key.String(), value)
}
