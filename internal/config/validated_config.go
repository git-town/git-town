package config

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type ValidatedConfig struct {
	*UnvalidatedConfig
	Config configdomain.ValidatedConfig // the merged configuration data
}

func EmptyValidatedConfig() ValidatedConfig {
	return ValidatedConfig{} //exhaustruct:ignore
}

func (self *ValidatedConfig) Reload() {
	_, self.GlobalGitConfig, _ = self.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, self.LocalGitConfig, _ = self.GitConfig.LoadLocal(false)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	unvalidateConfig := configdomain.NewUnvalidatedConfig(self.ConfigFile, self.GlobalGitConfig, self.LocalGitConfig)
	self.UnvalidatedConfig = &UnvalidatedConfig{
		Config:          NewMutable(&unvalidateConfig),
		ConfigFile:      self.ConfigFile,
		DryRun:          self.DryRun,
		GitConfig:       self.GitConfig,
		GitVersion:      self.GitVersion,
		GlobalGitConfig: self.GlobalGitConfig,
		LocalGitConfig:  self.LocalGitConfig,
	}
	self.Config = configdomain.ValidatedConfig{
		SharedConfig: unvalidateConfig.SharedConfig,
		GitUserEmail: self.Config.GitUserEmail,
		GitUserName:  self.Config.GitUserName,
		MainBranch:   self.Config.MainBranch,
	}
}
