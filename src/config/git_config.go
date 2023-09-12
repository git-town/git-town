package config

// GitConfig is an in-memory representation of the total Git configuration, global and local.
type GitConfig struct {
	Global GitConfigCache
	Local  GitConfigCache
}

func LoadGitConfig(runner runner) GitConfig {
	return GitConfig{
		Global: LoadGitConfigCache(runner, true),
		Local:  LoadGitConfigCache(runner, false),
	}
}

func (gc GitConfig) Clone() GitConfig {
	return GitConfig{
		Global: gc.Global.Clone(),
		Local:  gc.Local.Clone(),
	}
}
