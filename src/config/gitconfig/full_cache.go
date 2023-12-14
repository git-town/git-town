package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// FullCache caches all Git-based configuration types (global and local).
type FullCache struct {
	GlobalCache  SingleCache
	GlobalConfig configdomain.PartialConfig
	LocalCache   SingleCache
	LocalConfig  configdomain.PartialConfig
}

func LoadFullCache(git *Access) FullCache {
	return FullCache{
		GlobalCache: git.LoadCache(true),
		LocalCache:  git.LoadCache(false),
	}
}

func (self FullCache) Clone() FullCache {
	return FullCache{
		GlobalCache:  self.GlobalCache.Clone(),
		GlobalConfig: self.GlobalConfig,
		LocalCache:   self.LocalCache.Clone(),
		LocalConfig:  self.LocalConfig,
	}
}
