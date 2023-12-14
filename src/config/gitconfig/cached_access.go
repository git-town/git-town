package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// CachedAccess provides access to the local and global configuration data stored in Git metadata
// made efficient through an in-memory cache.
type CachedAccess struct {
	Access
	LocalGlobal
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewGit(gitConfig LocalGlobal, runner Runner) CachedAccess {
	return CachedAccess{
		LocalGlobal: gitConfig,
		Access: Access{
			Runner: runner,
		},
	}
}

func (self CachedAccess) GlobalConfigClone() Cache {
	return self.LocalGlobal.Global.Clone()
}

func (self CachedAccess) GlobalConfigValue(key configdomain.Key) string {
	return self.LocalGlobal.Global[key]
}

func (self CachedAccess) LocalConfigClone() Cache {
	return self.LocalGlobal.Local.Clone()
}

func (self CachedAccess) LocalConfigKeysMatching(pattern string) []configdomain.Key {
	return self.LocalGlobal.Local.KeysMatching(pattern)
}

func (self CachedAccess) LocalConfigValue(key configdomain.Key) string {
	return self.LocalGlobal.Local[key]
}

// LocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (self CachedAccess) LocalOrGlobalConfigValue(key configdomain.Key) string {
	local := self.LocalConfigValue(key)
	if local != "" {
		return local
	}
	return self.GlobalConfigValue(key)
}

// Reload refreshes the cached configuration information.
func (self *CachedAccess) Reload() {
	self.LocalGlobal = LoadLocalGlobal(&self.Access)
}

func (self *CachedAccess) RemoveGlobalConfigValue(key configdomain.Key) error {
	delete(self.LocalGlobal.Global, key)
	return self.Access.RemoveGlobalConfigValue(key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *CachedAccess) RemoveLocalConfigValue(key configdomain.Key) error {
	delete(self.LocalGlobal.Local, key)
	return self.Access.RemoveLocalConfigValue(key)
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *CachedAccess) SetGlobalConfigValue(key configdomain.Key, value string) error {
	self.LocalGlobal.Global[key] = value
	return self.Access.SetGlobalConfigValue(key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *CachedAccess) SetLocalConfigValue(key configdomain.Key, value string) error {
	self.LocalGlobal.Local[key] = value
	return self.Access.SetLocalConfigValue(key, value)
}
