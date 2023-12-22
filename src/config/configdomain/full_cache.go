package configdomain

// FullCache caches all Git-based configuration types (global and local).
type FullCache struct {
	GlobalCache  SingleCache
	GlobalConfig PartialConfig
	LocalCache   SingleCache
	LocalConfig  PartialConfig
}

func EmptyFullCache() FullCache {
	return FullCache{
		GlobalCache:  map[Key]string{},
		GlobalConfig: PartialConfig{}, //nolint:exhaustruct
		LocalCache:   map[Key]string{},
		LocalConfig:  PartialConfig{}, //nolint:exhaustruct
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
