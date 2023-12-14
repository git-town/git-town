package gitconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// GitConfig is an in-memory representation of the total Git configuration, global and local.
type GitConfig struct {
	GlobalCache  Cache
	GlobalConfig configdomain.PartialConfig
	LocalCache   Cache
	LocalConfig  configdomain.PartialConfig
}

func LoadGitConfig(runner Runner) GitConfig {
	globalConfig, globalCache := LoadGitConfigCache(runner, true)
	localConfig, localCache := LoadGitConfigCache(runner, false)
	return GitConfig{
		GlobalCache:  globalCache,
		GlobalConfig: globalConfig,
		LocalCache:   localCache,
		LocalConfig:  localConfig,
	}
}

func (self GitConfig) Clone() GitConfig {
	return GitConfig{
		GlobalCache:  self.GlobalCache.Clone(),
		GlobalConfig: self.GlobalConfig,
		LocalCache:   self.LocalCache.Clone(),
		LocalConfig:  self.LocalConfig,
	}
}
