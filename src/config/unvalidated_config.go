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
	Config          *configdomain.UnvalidatedConfig    // the merged configuration data
	ConfigFile      Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun          bool
	GitConfig       gitconfig.Access            // access to the Git configuration settings
	GlobalGitConfig configdomain.PartialConfig  // content of the global Git configuration
	LocalGitConfig  configdomain.PartialConfig  // content of the local Git configuration
	originURLCache  configdomain.OriginURLCache // TODO: remove if unused
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) (UnvalidatedConfig, stringslice.Collector) {
	config := configdomain.NewUnvalidatedConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	finalMessages := stringslice.NewCollector()
	return UnvalidatedConfig{
		Config:          config,
		ConfigFile:      args.ConfigFile,
		DryRun:          args.DryRun,
		GitConfig:       args.Access,
		GlobalGitConfig: args.GlobalConfig,
		LocalGitConfig:  args.LocalConfig,
		originURLCache:  configdomain.OriginURLCache{},
	}, finalMessages
}

// AddToContributionBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *UnvalidatedConfig) AddToContributionBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetContributionBranches(append(self.Config.ContributionBranches, branches...))
}

// AddToObservedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *UnvalidatedConfig) AddToObservedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetObservedBranches(append(self.Config.ObservedBranches, branches...))
}

// AddToParkedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *UnvalidatedConfig) AddToParkedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetParkedBranches(append(self.Config.ParkedBranches, branches...))
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *UnvalidatedConfig) OriginURL() Option[giturl.Parts] {
	text := self.OriginURLString()
	if text == "" {
		return None[giturl.Parts]()
	}
	return confighelpers.DetermineOriginURL(text, self.Config.HostingOriginHostname, self.originURLCache)
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
	self.Config.ContributionBranches = slice.Remove(self.Config.ContributionBranches, branch)
	return self.SetContributionBranches(self.Config.ContributionBranches)
}

// RemoveFromParkedBranches removes the given branch as a perennial branch.
func (self *UnvalidatedConfig) RemoveFromParkedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ParkedBranches = slice.Remove(self.Config.ParkedBranches, branch)
	return self.SetParkedBranches(self.Config.ParkedBranches)
}

// RemoveFromObservedBranches removes the given branch as a perennial branch.
func (self *UnvalidatedConfig) RemoveFromObservedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ObservedBranches = slice.Remove(self.Config.ObservedBranches, branch)
	return self.SetObservedBranches(self.Config.ObservedBranches)
}

func (self *UnvalidatedConfig) RemoveMainBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyMainBranch)
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (self *UnvalidatedConfig) RemoveOutdatedConfiguration(localBranches gitdomain.LocalBranchNames) error {
	for child, parent := range self.Config.Lineage {
		hasChildBranch := localBranches.Contains(child)
		hasParentBranch := localBranches.Contains(parent)
		if !hasChildBranch || !hasParentBranch {
			self.RemoveParent(child)
		}
	}
	return nil
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *UnvalidatedConfig) RemoveParent(branch gitdomain.LocalBranchName) {
	if self.LocalGitConfig.Lineage != nil {
		self.LocalGitConfig.Lineage.RemoveBranch(branch)
	}
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

func (self *UnvalidatedConfig) RemoveSyncBeforeShip() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncBeforeShip)
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
	self.Config.ContributionBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyContributionBranches, branches.Join(" "))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.Config.MainBranch = Some(branch)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, branch.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *UnvalidatedConfig) SetObservedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ObservedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyObservedBranches, branches.Join(" "))
}

// SetOffline updates whether Git Town is in offline mode.
func (self *UnvalidatedConfig) SetOffline(value configdomain.Offline) error {
	self.Config.Offline = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Config.Lineage[branch] = parentBranch
	return self.GitConfig.SetLocalConfigValue(gitconfig.NewParentKey(branch), parentBranch.String())
}

// SetObservedBranches marks the given branches as perennial branches.
func (self *UnvalidatedConfig) SetParkedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ParkedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyParkedBranches, branches.Join(" "))
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *UnvalidatedConfig) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.PerennialBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, branches.Join(" "))
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *UnvalidatedConfig) SetPerennialRegexLocally(value configdomain.PerennialRegex) error {
	self.Config.PerennialRegex = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialRegex, value.String())
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *UnvalidatedConfig) SetPushHookLocally(value configdomain.PushHook) error {
	self.Config.PushHook = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetPushNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *UnvalidatedConfig) SetPushNewBranches(value configdomain.PushNewBranches, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.Config.PushNewBranches = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushNewBranches, setting)
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushNewBranches, setting)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *UnvalidatedConfig) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, global bool) error {
	self.Config.ShipDeleteTrackingBranch = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *UnvalidatedConfig) SetSyncBeforeShip(value configdomain.SyncBeforeShip, global bool) error {
	self.Config.SyncBeforeShip = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
}

func (self *UnvalidatedConfig) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.Config.SyncFeatureStrategy = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *UnvalidatedConfig) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.Config.SyncPerennialStrategy = strategy
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *UnvalidatedConfig) SetSyncUpstream(value configdomain.SyncUpstream, global bool) error {
	self.Config.SyncUpstream = value
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
