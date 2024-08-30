package config

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/slice"
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

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *ValidatedConfig) AddToPerennialBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetPerennialBranches(append(self.Config.PerennialBranches, branches...))
}

// Author provides the locally Git configured user.
func (self *ValidatedConfig) Author() gitdomain.Author {
	email := self.Config.GitUserEmail
	name := self.Config.GitUserName
	return gitdomain.Author(fmt.Sprintf("%s <%s>", name, email))
}

func (self *ValidatedConfig) Reload() {
	_, self.GlobalGitConfig, _ = self.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, self.LocalGitConfig, _ = self.GitConfig.LoadLocal(false)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	unvalidateConfig := configdomain.NewUnvalidatedConfig(self.ConfigFile, self.GlobalGitConfig, self.LocalGitConfig)
	self.Config = configdomain.ValidatedConfig{
		UnvalidatedConfig: &unvalidateConfig,
		GitUserEmail:      self.Config.GitUserEmail,
		GitUserName:       self.Config.GitUserName,
		MainBranch:        self.Config.MainBranch,
	}
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *ValidatedConfig) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	self.Config.PerennialBranches = slice.Remove(self.Config.PerennialBranches, branch)
	return self.SetPerennialBranches(self.Config.PerennialBranches)
}
