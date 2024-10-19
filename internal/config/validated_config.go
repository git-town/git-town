package config

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type ValidatedConfig struct {
	NormalConfig    NormalConfig
	ValidatedConfig configdomain.ValidatedConfig
}

func EmptyValidatedConfig() ValidatedConfig {
	return ValidatedConfig{} //exhaustruct:ignore
}

func (self *ValidatedConfig) BranchType(branch gitdomain.LocalBranchName) configdomain.BranchType {
	if self.ValidatedConfig.IsMainBranch(branch) {
		return configdomain.BranchTypeMainBranch
	}
	return self.NormalConfig.PartialBranchType(branch)
}

func (self *ValidatedConfig) Reload() {
	_, globalGitConfig, _ := self.NormalConfig.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, localGitConfig, _ := self.NormalConfig.GitConfig.LoadLocal(false)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	self.ValidatedConfig = configdomain.NewValidatedConfig(self.NormalConfig.ConfigFile, self.NormalConfig.GlobalGitConfig, self.NormalConfig.LocalGitConfig, self.ValidatedConfig)
	unvalidatedConfig, _ := NewUnvalidatedConfig(NewUnvalidatedConfigArgs{
		Access:       self.NormalConfig.GitConfig,
		ConfigFile:   self.NormalConfig.ConfigFile,
		DryRun:       self.NormalConfig.DryRun,
		GitVersion:   self.NormalConfig.GitVersion,
		GlobalConfig: globalGitConfig,
		LocalConfig:  localGitConfig,
	})
	self.ValidatedConfig = unvalidatedConfig.ToValidatedConfig()
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *ValidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.ValidatedConfig.MainBranch = branch
	return self.NormalConfig.GitConfig.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
}
