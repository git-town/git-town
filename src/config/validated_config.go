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
	UnvalidatedConfig
	Config configdomain.ValidatedConfig // the merged configuration data
}

// AddToContributionBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *ValidatedConfig) AddToContributionBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetContributionBranches(append(self.Config.ContributionBranches, branches...))
}

// AddToObservedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *ValidatedConfig) AddToObservedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetObservedBranches(append(self.Config.ObservedBranches, branches...))
}

// AddToParkedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *ValidatedConfig) AddToParkedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetParkedBranches(append(self.Config.ParkedBranches, branches...))
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
		MainBranch:        self.Config.MainBranch,
	}
}

// RemoveFromContributionBranches removes the given branch as a perennial branch.
func (self *ValidatedConfig) RemoveFromContributionBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ContributionBranches = slice.Remove(self.Config.ContributionBranches, branch)
	return self.SetContributionBranches(self.Config.ContributionBranches)
}

// RemoveFromObservedBranches removes the given branch as a perennial branch.
func (self *ValidatedConfig) RemoveFromObservedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ObservedBranches = slice.Remove(self.Config.ObservedBranches, branch)
	return self.SetObservedBranches(self.Config.ObservedBranches)
}

// RemoveFromParkedBranches removes the given branch as a perennial branch.
func (self *ValidatedConfig) RemoveFromParkedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ParkedBranches = slice.Remove(self.Config.ParkedBranches, branch)
	return self.SetParkedBranches(self.Config.ParkedBranches)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *ValidatedConfig) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	self.Config.PerennialBranches = slice.Remove(self.Config.PerennialBranches, branch)
	return self.SetPerennialBranches(self.Config.PerennialBranches)
}

func (self *ValidatedConfig) RemoveMainBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyMainBranch)
}

func (self *ValidatedConfig) RemovePerennialBranches() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPerennialBranches)
}

func (self *ValidatedConfig) RemovePerennialRegex() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPerennialRegex)
}

func (self *ValidatedConfig) RemovePushHook() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPushHook)
}

func (self *ValidatedConfig) RemovePushNewBranches() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPushNewBranches)
}

func (self *ValidatedConfig) RemoveShipDeleteTrackingBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch)
}

func (self *ValidatedConfig) RemoveSyncBeforeShip() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncBeforeShip)
}

func (self *ValidatedConfig) RemoveSyncFeatureStrategy() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncFeatureStrategy)
}

func (self *ValidatedConfig) RemoveSyncPerennialStrategy() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncPerennialStrategy)
}

func (self *ValidatedConfig) RemoveSyncUpstream() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncUpstream)
}

// SetObservedBranches marks the given branches as observed branches.
func (self *ValidatedConfig) SetContributionBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ContributionBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyContributionBranches, branches.Join(" "))
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *ValidatedConfig) SetObservedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ObservedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyObservedBranches, branches.Join(" "))
}

// SetOriginHostname marks the given branch as the main branch
// in the Git Town configuration.
func (self *ValidatedConfig) SetOriginHostname(hostName configdomain.HostingOriginHostname) error {
	self.Config.HostingOriginHostname = Some(hostName)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyHostingOriginHostname, hostName.String())
}

// SetObservedBranches marks the given branches as perennial branches.
func (self *ValidatedConfig) SetParkedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ParkedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyParkedBranches, branches.Join(" "))
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *ValidatedConfig) SetPerennialRegexLocally(value configdomain.PerennialRegex) error {
	self.Config.PerennialRegex = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialRegex, value.String())
}

// SetPushHook updates the configured push-hook strategy.
func (self *ValidatedConfig) SetPushHookGlobally(value configdomain.PushHook) error {
	self.Config.PushHook = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(value.Bool()))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *ValidatedConfig) SetPushHookLocally(value configdomain.PushHook) error {
	self.Config.PushHook = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetPushNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *ValidatedConfig) SetPushNewBranches(value configdomain.PushNewBranches, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.Config.PushNewBranches = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushNewBranches, setting)
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushNewBranches, setting)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *ValidatedConfig) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, global bool) error {
	self.Config.ShipDeleteTrackingBranch = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *ValidatedConfig) SetSyncBeforeShip(value configdomain.SyncBeforeShip, global bool) error {
	self.Config.SyncBeforeShip = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
}

func (self *ValidatedConfig) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.Config.SyncFeatureStrategy = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

func (self *ValidatedConfig) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	self.Config.SyncFeatureStrategy = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *ValidatedConfig) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.Config.SyncPerennialStrategy = strategy
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *ValidatedConfig) SetSyncUpstream(value configdomain.SyncUpstream, global bool) error {
	self.Config.SyncUpstream = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}
