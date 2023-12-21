package gitconfig

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

// FullCache caches all Git-based configuration types (global and local).
type FullCache struct {
	GlobalCache  SingleCache
	GlobalConfig configdomain.PartialConfig
	LocalCache   SingleCache
	LocalConfig  configdomain.PartialConfig
}

func EmptyFullCache() FullCache {
	return FullCache{
		GlobalCache:  map[configdomain.Key]string{},
		GlobalConfig: configdomain.PartialConfig{}, //nolint:exhaustruct
		LocalCache:   map[configdomain.Key]string{},
		LocalConfig:  configdomain.PartialConfig{}, //nolint:exhaustruct
	}
}

func LoadFullCache(access *Access) (FullCache, error) {
	globalCache, globalConfig, err := access.LoadCache(true, access.RemoveGlobalConfigValue)
	if err != nil {
		return EmptyFullCache(), err
	}
	localCache, localConfig, err := access.LoadCache(false, access.RemoveLocalConfigValue)
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
