package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// Cached manages configuration data stored in Cached metadata.
// Supports configuration in the local repo and the global Cached configuration.
type Cached struct {
	Git
	LocalGlobal
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(gitConfig LocalGlobal, runner Runner) Cached {
	return Cached{
		LocalGlobal: gitConfig,
		Git: Git{
			Runner: runner,
		},
	}
}

func (self Cached) GlobalConfigClone() Cache {
	return self.LocalGlobal.Global.Clone()
}

func (self Cached) GlobalConfigValue(key configdomain.Key) string {
	return self.LocalGlobal.Global[key]
}

func (self Cached) LocalConfigClone() Cache {
	return self.LocalGlobal.Local.Clone()
}

func (self Cached) LocalConfigKeysMatching(pattern string) []configdomain.Key {
	return self.LocalGlobal.Local.KeysMatching(pattern)
}

func (self Cached) LocalConfigValue(key configdomain.Key) string {
	return self.LocalGlobal.Local[key]
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
	self.LocalGlobal = LoadLocalGlobal(self.Runner)
}

func (self *Cached) RemoveGlobalConfigValue(key configdomain.Key) error {
	delete(self.LocalGlobal.Global, key)
	return self.Git.RemoveGlobalConfigValue(key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *Cached) RemoveLocalConfigValue(key configdomain.Key) error {
	delete(self.LocalGlobal.Local, key)
	return self.Git.RemoveLocalConfigValue(key)
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *Cached) SetGlobalConfigValue(key configdomain.Key, value string) error {
	self.LocalGlobal.Global[key] = value
	return self.Git.SetGlobalConfigValue(key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *Cached) SetLocalConfigValue(key configdomain.Key, value string) error {
	self.LocalGlobal.Local[key] = value
	return self.Git.SetLocalConfigValue(key, value)
}
