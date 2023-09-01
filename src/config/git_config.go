package config

// GitConfig represents the total Git configuration, global and local.
type GitConfig struct {
	Global GitConfigCache
	Local  GitConfigCache
}

func LoadGitConfig(querier querier) GitConfig {
	return GitConfig{
		Global: LoadGitConfigCache(querier, true),
		Local:  LoadGitConfigCache(querier, false),
	}
}
