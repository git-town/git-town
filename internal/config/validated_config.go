package config

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
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

// TODO: the ValidatedConfig should not reload itself from disk since it is just a part of the overall config now.
// Instead, create a static function that loads the config from disk and returns NormalConfig and ValidatedConfig.
func (self *ValidatedConfig) Reload() {
	_, self.GlobalGitConfig, _ = self.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	self.Config = configdomain.ValidatedConfig{
		GitUserEmail: self.Config.GitUserEmail,
		GitUserName:  self.Config.GitUserName,
		MainBranch:   self.Config.MainBranch,
	}
}
