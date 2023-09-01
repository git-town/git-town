package config

// GitConfig is an in-memory representation of the total Git configuration, global and local.
type GitConfig struct {
	Global GitConfigCache
	Local  GitConfigCache
}

func LoadGitConfig(runner runner) GitConfig {
	return GitConfig{
		Local:  LoadGitConfigCache(runner, false),
		Global: LoadGitConfigCache(runner, true),
	}
}
