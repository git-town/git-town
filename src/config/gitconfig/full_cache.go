package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// FullCache caches all Git-based configuration types (global and local).
type FullCache struct {
	GlobalCache  SingleCache
	GlobalConfig configdomain.PartialGitConfig
	LocalCache   SingleCache
	LocalConfig  configdomain.PartialGitConfig
}

func EmptyFullCache() FullCache {
	return FullCache{
		GlobalCache:  map[configdomain.Key]string{},
		GlobalConfig: configdomain.PartialGitConfig{}, //nolint:exhaustruct
		LocalCache:   map[configdomain.Key]string{},
		LocalConfig:  configdomain.PartialGitConfig{}, //nolint:exhaustruct
	}
}

func LoadFullCache(access *Access) (FullCache, error) {
	globalCache, globalConfig, err := access.LoadCache(true)
	if err != nil {
		return EmptyFullCache(), err
	}
	localCache, localConfig, err := access.LoadCache(false)
	return FullCache{
		GlobalCache:  globalCache,
		GlobalConfig: globalConfig,
		LocalCache:   localCache,
		LocalConfig:  localConfig,
	}, err
}

func (self FullCache) Clone() FullCache {
	return FullCache{
		GlobalCache:  self.GlobalCache.Clone(),
		GlobalConfig: self.GlobalConfig,
		LocalCache:   self.LocalCache.Clone(),
		LocalConfig:  self.LocalConfig,
	}
}
