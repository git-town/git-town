package gitconfig

// LocalGlobal is an in-memory representation of the entire Git configuration, global and local.
type LocalGlobal struct {
	Global Cache
	Local  Cache
}

func LoadLocalGlobal(git *Git) LocalGlobal {
	return LocalGlobal{
		Global: git.LoadCache(true),
		Local:  git.LoadCache(false),
	}
}

func (self LocalGlobal) Clone() LocalGlobal {
	return LocalGlobal{
		Global: self.Global.Clone(),
		Local:  self.Local.Clone(),
	}
}
