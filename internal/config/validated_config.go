package config

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type ValidatedConfig struct {
	NormalConfig    NormalConfig
	ValidatedConfig configdomain.ValidatedConfig // the merged configuration data
}

func EmptyValidatedConfig() ValidatedConfig {
	return ValidatedConfig{} //exhaustruct:ignore
}

func (self *ValidatedConfig) Reload() {
	_, self.GlobalGitConfig, _ = self.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, self.LocalGitConfig, _ = self.GitConfig.LoadLocal(false)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	validateConfig := configdomain.NewValidatedConfig(self.ConfigFile, self.GlobalGitConfig, self.LocalGitConfig, self.ValidatedConfig)
	self.ValidatedConfig = configdomain.ValidatedConfig{
		NormalConfig: validateConfig.NormalConfig,
		GitUserEmail: self.ValidatedConfig.GitUserEmail,
		GitUserName:  self.ValidatedConfig.GitUserName,
		MainBranch:   self.ValidatedConfig.MainBranch,
	}
}
