package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// CachedAccess provides access to the local and global configuration data stored in Git metadata
// made efficient through an in-memory cache.
type CachedAccess struct {
	Access
	FullCache
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewCachedAccess(fullCache FullCache, runner Runner) CachedAccess {
	return CachedAccess{
		FullCache: fullCache,
		Access: Access{
			Runner: runner,
		},
	}
}

func (self CachedAccess) GlobalConfigValue(key configdomain.Key) string {
	return self.GlobalCache[key]
}

func (self CachedAccess) LocalConfigKeysMatching(pattern string) []configdomain.Key {
	return self.FullCache.LocalCache.KeysMatching(pattern)
}

func (self CachedAccess) LocalConfigValue(key configdomain.Key) string {
	return self.FullCache.LocalCache[key]
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
	self.FullCache = LoadFullCache(&self.Access)
}

func (self *CachedAccess) RemoveGlobalConfigValue(key configdomain.Key) error {
	delete(self.GlobalCache, key)
	return self.Access.RemoveGlobalConfigValue(key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (self *CachedAccess) RemoveLocalConfigValue(key configdomain.Key) error {
	delete(self.FullCache.LocalCache, key)
	return self.Access.RemoveLocalConfigValue(key)
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *CachedAccess) SetGlobalConfigValue(key configdomain.Key, value string) error {
	self.GlobalCache[key] = value
	return self.Access.SetGlobalConfigValue(key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *CachedAccess) SetLocalConfigValue(key configdomain.Key, value string) error {
	self.FullCache.LocalCache[key] = value
	return self.Access.SetLocalConfigValue(key, value)
}
