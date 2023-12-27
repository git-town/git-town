package configdomain

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

// Reload refreshes the cached configuration information.
func (self *CachedAccess) Reload() {
	self.FullCache, _ = LoadFullCache(&self.Access)
}
