package configdomain

// FullCache caches all Git-based configuration types (global and local).
type FullCache struct {
	GlobalCache SingleCache
	LocalCache  SingleCache
}

func EmptyFullCache() FullCache {
	return FullCache{
		GlobalCache: map[Key]string{},
		LocalCache:  map[Key]string{},
	}
}

func LoadFullCache(access *Access) (FullCache, error) {
	globalCache, globalConfig, err := access.LoadCache(true)
	if err != nil {
		return EmptyFullCache(), err
	}
	localCache, localConfig, err := access.LoadCache(false)
	return FullCache{
		GlobalCache: globalCache,
		LocalCache:  localCache,
	}, err
}

func (self FullCache) Clone() FullCache {
	return FullCache{
		GlobalCache: self.GlobalCache.Clone(),
		LocalCache:  self.LocalCache.Clone(),
	}
}
