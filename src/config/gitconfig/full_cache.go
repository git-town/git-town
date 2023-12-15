package gitconfig

// FullCache caches all Git-based configuration types (global and local).
type FullCache struct {
	GlobalCache SingleCache
	LocalCache  SingleCache
}

func LoadFullCache(git *Access) FullCache {
	return FullCache{
		GlobalCache: git.LoadCache(true),
		LocalCache:  git.LoadCache(false),
	}
}

func (self FullCache) Clone() FullCache {
	return FullCache{
		GlobalCache: self.GlobalCache.Clone(),
		LocalCache:  self.LocalCache.Clone(),
	}
}
