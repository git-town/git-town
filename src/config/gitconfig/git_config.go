package gitconfig

// GitConfig is an in-memory representation of the total Git configuration, global and local.
type GitConfig struct {
	Global Cache
	Local  Cache
}

func LoadGitConfig(runner Runner) GitConfig {
	return GitConfig{
		Global: LoadGitConfigCache(runner, true),
		Local:  LoadGitConfigCache(runner, false),
	}
}

func (self GitConfig) Clone() GitConfig {
	return GitConfig{
		Global: self.Global.Clone(),
		Local:  self.Local.Clone(),
	}
}
