package config

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
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
	_, self.NormalConfig.GlobalGitConfig, _ = self.NormalConfig.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, self.NormalConfig.LocalGitConfig, _ = self.NormalConfig.GitConfig.LoadLocal(false)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	validateConfig := configdomain.NewValidatedConfig(self.NormalConfig.ConfigFile, self.NormalConfig.GlobalGitConfig, self.NormalConfig.LocalGitConfig, self.ValidatedConfig)
	self.ValidatedConfig = configdomain.ValidatedConfig{
		NormalConfig: validateConfig.NormalConfig,
		GitUserEmail: self.ValidatedConfig.GitUserEmail,
		GitUserName:  self.ValidatedConfig.GitUserName,
		MainBranch:   self.ValidatedConfig.MainBranch,
	}
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *ValidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.ValidatedConfig.MainBranch = branch
	return self.NormalConfig.GitConfig.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
}
