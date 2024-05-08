package config

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
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
	self.Config = configdomain.ValidatedConfig{
		UnvalidatedConfig: configdomain.NewUnvalidatedConfig(self.ConfigFile, self.GlobalGitConfig, self.LocalGitConfig),
		GitUserEmail:      self.Config.GitUserEmail,
		GitUserName:       self.Config.GitUserName,
		MainBranch:        self.Config.MainBranch,
	}
}

// RemoveFromContributionBranches removes the given branch as a perennial branch.
func (self *ValidatedConfig) RemoveFromContributionBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ContributionBranches = slice.Remove(self.Config.ContributionBranches, branch)
	return self.SetContributionBranches(self.Config.ContributionBranches)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *ValidatedConfig) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	self.Config.PerennialBranches = slice.Remove(self.Config.PerennialBranches, branch)
	return self.SetPerennialBranches(self.Config.PerennialBranches)
}

// SetOriginHostname marks the given branch as the main branch
// in the Git Town configuration.
func (self *ValidatedConfig) SetOriginHostname(hostName configdomain.HostingOriginHostname) error {
	self.Config.HostingOriginHostname = Some(hostName)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyHostingOriginHostname, hostName.String())
}

// SetPushHook updates the configured push-hook strategy.
func (self *ValidatedConfig) SetPushHookGlobally(value configdomain.PushHook) error {
	self.Config.PushHook = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(value.Bool()))
}

func (self *ValidatedConfig) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	self.Config.SyncFeatureStrategy = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}
