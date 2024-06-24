package config

import (
	"strconv"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/confighelpers"
	"github.com/git-town/git-town/v14/src/config/envconfig"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
)

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

type UnvalidatedConfig struct {
	Config          Mutable[configdomain.UnvalidatedConfig] // the merged configuration data
	ConfigFile      Option[configdomain.PartialConfig]      // content of git-town.toml, nil = no config file exists
	DryRun          bool
	GitConfig       gitconfig.Access           // access to the Git configuration settings
	GlobalGitConfig configdomain.PartialConfig // content of the global Git configuration
	LocalGitConfig  configdomain.PartialConfig // content of the local Git configuration
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) (UnvalidatedConfig, stringslice.Collector) {
	config := configdomain.NewUnvalidatedConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	finalMessages := stringslice.NewCollector()
	return UnvalidatedConfig{
		Config:          NewMutable(&config),
		ConfigFile:      args.ConfigFile,
		DryRun:          args.DryRun,
		GitConfig:       args.Access,
		GlobalGitConfig: args.GlobalConfig,
		LocalGitConfig:  args.LocalConfig,
	}, finalMessages
}

// AddToContributionBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *UnvalidatedConfig) AddToContributionBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetContributionBranches(append(self.Config.Value.ContributionBranches, branches...))
}

// AddToObservedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *UnvalidatedConfig) AddToObservedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetObservedBranches(append(self.Config.Value.ObservedBranches, branches...))
}

// AddToParkedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *UnvalidatedConfig) AddToParkedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetParkedBranches(append(self.Config.Value.ParkedBranches, branches...))
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *UnvalidatedConfig) OriginURL() Option[giturl.Parts] {
	text := self.OriginURLString()
	if text == "" {
		return None[giturl.Parts]()
	}
	return confighelpers.DetermineOriginURL(text, self.Config.Value.HostingOriginHostname)
}

// OriginURLString provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *UnvalidatedConfig) OriginURLString() string {
	remoteOverride := envconfig.OriginURLOverride()
	if remoteOverride != "" {
		return remoteOverride
	}
	return self.GitConfig.OriginRemote()
}

// RemoveFromContributionBranches removes the given branch as a perennial branch.
func (self *UnvalidatedConfig) RemoveFromContributionBranches(branch gitdomain.LocalBranchName) error {
	self.Config.Value.ContributionBranches = slice.Remove(self.Config.Value.ContributionBranches, branch)
	return self.SetContributionBranches(self.Config.Value.ContributionBranches)
}

// RemoveFromObservedBranches removes the given branch as a perennial branch.
func (self *UnvalidatedConfig) RemoveFromObservedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.Value.ObservedBranches = slice.Remove(self.Config.Value.ObservedBranches, branch)
	return self.SetObservedBranches(self.Config.Value.ObservedBranches)
}

// RemoveFromParkedBranches removes the given branch as a perennial branch.
func (self *UnvalidatedConfig) RemoveFromParkedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.Value.ParkedBranches = slice.Remove(self.Config.Value.ParkedBranches, branch)
	return self.SetParkedBranches(self.Config.Value.ParkedBranches)
}

func (self *UnvalidatedConfig) RemoveMainBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyMainBranch)
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (self *UnvalidatedConfig) RemoveOutdatedConfiguration(localBranches gitdomain.LocalBranchNames) error {
	for _, entry := range self.Config.Value.Lineage.Entries() {
		hasChildBranch := localBranches.Contains(entry.Child)
		hasParentBranch := localBranches.Contains(entry.Parent)
		if !hasChildBranch || !hasParentBranch {
			self.RemoveParent(entry.Child)
		}
	}
	return nil
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *UnvalidatedConfig) RemoveParent(branch gitdomain.LocalBranchName) {
	self.LocalGitConfig.Lineage.RemoveBranch(branch)
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.NewParentKey(branch))
}

func (self *UnvalidatedConfig) RemovePerennialBranches() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPerennialBranches)
}

func (self *UnvalidatedConfig) RemovePerennialRegex() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPerennialRegex)
}

func (self *UnvalidatedConfig) RemovePushHook() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPushHook)
}

func (self *UnvalidatedConfig) RemovePushNewBranches() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPushNewBranches)
}

func (self *UnvalidatedConfig) RemoveShipDeleteTrackingBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch)
}

func (self *UnvalidatedConfig) RemoveSyncFeatureStrategy() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncFeatureStrategy)
}

func (self *UnvalidatedConfig) RemoveSyncPerennialStrategy() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncPerennialStrategy)
}

func (self *UnvalidatedConfig) RemoveSyncUpstream() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncUpstream)
}

// SetObservedBranches marks the given branches as observed branches.
func (self *UnvalidatedConfig) SetContributionBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.Value.ContributionBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyContributionBranches, branches.Join(" "))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.Config.Value.MainBranch = Some(branch)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, branch.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *UnvalidatedConfig) SetObservedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.Value.ObservedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyObservedBranches, branches.Join(" "))
}

// SetOffline updates whether Git Town is in offline mode.
func (self *UnvalidatedConfig) SetOffline(value configdomain.Offline) error {
	self.Config.Value.Offline = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Config.Value.Lineage.Add(branch, parentBranch)
	return self.GitConfig.SetLocalConfigValue(gitconfig.NewParentKey(branch), parentBranch.String())
}

// SetObservedBranches marks the given branches as perennial branches.
func (self *UnvalidatedConfig) SetParkedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.Value.ParkedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyParkedBranches, branches.Join(" "))
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *UnvalidatedConfig) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.Value.PerennialBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, branches.Join(" "))
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *UnvalidatedConfig) SetPerennialRegexLocally(value configdomain.PerennialRegex) error {
	self.Config.Value.PerennialRegex = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialRegex, value.String())
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *UnvalidatedConfig) SetPushHookLocally(value configdomain.PushHook) error {
	self.Config.Value.PushHook = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetPushNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *UnvalidatedConfig) SetPushNewBranches(value configdomain.PushNewBranches, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.Config.Value.PushNewBranches = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushNewBranches, setting)
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushNewBranches, setting)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *UnvalidatedConfig) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, global bool) error {
	self.Config.Value.ShipDeleteTrackingBranch = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *UnvalidatedConfig) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.Config.Value.SyncFeatureStrategy = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *UnvalidatedConfig) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.Config.Value.SyncPerennialStrategy = strategy
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *UnvalidatedConfig) SetSyncUpstream(value configdomain.SyncUpstream, global bool) error {
	self.Config.Value.SyncUpstream = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}

type NewUnvalidatedConfigArgs struct {
	Access       gitconfig.Access
	ConfigFile   Option[configdomain.PartialConfig]
	DryRun       bool
	GlobalConfig configdomain.PartialConfig
	LocalConfig  configdomain.PartialConfig
}
