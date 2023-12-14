package gitconfig

// FullCache caches all Git-based configuration types (global and local).
type FullCache struct {
	Global SingleCache
	Local  SingleCache
}

func LoadFullCache(git *Access) FullCache {
	return FullCache{
		Global: git.LoadCache(true),
		Local:  git.LoadCache(false),
	}
}

func (self FullCache) Clone() FullCache {
	return FullCache{
		Global: self.Global.Clone(),
		Local:  self.Local.Clone(),
	}
}
