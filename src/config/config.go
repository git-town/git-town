package config

import (
	"fmt"
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

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type Config struct {
	Config          configdomain.FullConfig            // the merged configuration data
	ConfigFile      Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun          bool
	GitConfig       gitconfig.Access           // access to the Git configuration settings
	GlobalGitConfig configdomain.PartialConfig // content of the global Git configuration
	LocalGitConfig  configdomain.PartialConfig // content of the local Git configuration
	originURLCache  configdomain.OriginURLCache
}

// AddToContributionBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToContributionBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetContributionBranches(append(self.Config.ContributionBranches, branches...))
}

// AddToObservedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToObservedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetObservedBranches(append(self.Config.ObservedBranches, branches...))
}

// AddToParkedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToParkedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetParkedBranches(append(self.Config.ParkedBranches, branches...))
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToPerennialBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetPerennialBranches(append(self.Config.PerennialBranches, branches...))
}

// Author provides the locally Git configured user.
func (self *Config) Author() gitdomain.Author {
	email := self.Config.GitUserEmail
	name := self.Config.GitUserName
	return gitdomain.Author(fmt.Sprintf("%s <%s>", name, email))
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *Config) OriginURL() Option[giturl.Parts] {
	text := self.OriginURLString()
	if text == "" {
		return None[giturl.Parts]()
	}
	return confighelpers.DetermineOriginURL(text, self.Config.HostingOriginHostname, self.originURLCache)
}

// OriginURLString provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *Config) OriginURLString() string {
	remoteOverride := envconfig.OriginURLOverride()
	if remoteOverride != "" {
		return remoteOverride
	}
	return self.GitConfig.OriginRemote()
}

func (self *Config) Reload() {
	_, self.GlobalGitConfig, _ = self.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, self.LocalGitConfig, _ = self.GitConfig.LoadLocal(false)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	self.Config = configdomain.NewFullConfig(self.ConfigFile, self.GlobalGitConfig, self.LocalGitConfig)
}

// RemoveFromContributionBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromContributionBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ContributionBranches = slice.Remove(self.Config.ContributionBranches, branch)
	return self.SetContributionBranches(self.Config.ContributionBranches)
}

// RemoveFromObservedBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromObservedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ObservedBranches = slice.Remove(self.Config.ObservedBranches, branch)
	return self.SetObservedBranches(self.Config.ObservedBranches)
}

// RemoveFromParkedBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromParkedBranches(branch gitdomain.LocalBranchName) error {
	self.Config.ParkedBranches = slice.Remove(self.Config.ParkedBranches, branch)
	return self.SetParkedBranches(self.Config.ParkedBranches)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	self.Config.PerennialBranches = slice.Remove(self.Config.PerennialBranches, branch)
	return self.SetPerennialBranches(self.Config.PerennialBranches)
}

func (self *Config) RemoveMainBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyMainBranch)
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (self *Config) RemoveOutdatedConfiguration(localBranches gitdomain.LocalBranchNames) error {
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
func (self *Config) RemoveParent(branch gitdomain.LocalBranchName) {
	if self.LocalGitConfig.Lineage != nil {
		self.LocalGitConfig.Lineage.RemoveBranch(branch)
	}
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.NewParentKey(branch))
}

func (self *Config) RemovePerennialBranches() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPerennialBranches)
}

func (self *Config) RemovePerennialRegex() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPerennialRegex)
}

func (self *Config) RemovePushHook() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPushHook)
}

func (self *Config) RemovePushNewBranches() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyPushNewBranches)
}

func (self *Config) RemoveShipDeleteTrackingBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch)
}

func (self *Config) RemoveSyncBeforeShip() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncBeforeShip)
}

func (self *Config) RemoveSyncFeatureStrategy() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncFeatureStrategy)
}

func (self *Config) RemoveSyncPerennialStrategy() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncPerennialStrategy)
}

func (self *Config) RemoveSyncUpstream() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeySyncUpstream)
}

// SetObservedBranches marks the given branches as observed branches.
func (self *Config) SetContributionBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ContributionBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyContributionBranches, branches.Join(" "))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *Config) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.Config.MainBranch = branch
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, branch.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *Config) SetObservedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ObservedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyObservedBranches, branches.Join(" "))
}

// SetOffline updates whether Git Town is in offline mode.
func (self *Config) SetOffline(value configdomain.Offline) error {
	self.Config.Offline = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyOffline, value.String())
}

// SetOriginHostname marks the given branch as the main branch
// in the Git Town configuration.
func (self *Config) SetOriginHostname(hostName configdomain.HostingOriginHostname) error {
	self.Config.HostingOriginHostname = Some(hostName)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyHostingOriginHostname, hostName.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *Config) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Config.Lineage[branch] = parentBranch
	return self.GitConfig.SetLocalConfigValue(gitconfig.NewParentKey(branch), parentBranch.String())
}

// SetObservedBranches marks the given branches as perennial branches.
func (self *Config) SetParkedBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.ParkedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyParkedBranches, branches.Join(" "))
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *Config) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.PerennialBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, branches.Join(" "))
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *Config) SetPerennialRegexLocally(value configdomain.PerennialRegex) error {
	self.Config.PerennialRegex = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialRegex, value.String())
}

// SetPushHook updates the configured push-hook strategy.
func (self *Config) SetPushHookGlobally(value configdomain.PushHook) error {
	self.Config.PushHook = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(value.Bool()))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *Config) SetPushHookLocally(value configdomain.PushHook) error {
	self.Config.PushHook = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetPushNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *Config) SetPushNewBranches(value configdomain.PushNewBranches, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.Config.PushNewBranches = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushNewBranches, setting)
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushNewBranches, setting)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *Config) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, global bool) error {
	self.Config.ShipDeleteTrackingBranch = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *Config) SetSyncBeforeShip(value configdomain.SyncBeforeShip, global bool) error {
	self.Config.SyncBeforeShip = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
}

func (self *Config) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.Config.SyncFeatureStrategy = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

func (self *Config) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	self.Config.SyncFeatureStrategy = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *Config) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.Config.SyncPerennialStrategy = strategy
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *Config) SetSyncUpstream(value configdomain.SyncUpstream, global bool) error {
	self.Config.SyncUpstream = value
	if global {
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}

func NewConfig(args NewConfigArgs) (Config, *stringslice.Collector) {
	config := configdomain.NewFullConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	configAccess := gitconfig.Access{Runner: args.Runner}
	finalMessages := stringslice.NewCollector()
	return Config{
		Config:          config,
		ConfigFile:      args.ConfigFile,
		DryRun:          args.DryRun,
		GitConfig:       configAccess,
		GlobalGitConfig: args.GlobalConfig,
		LocalGitConfig:  args.LocalConfig,
		originURLCache:  configdomain.OriginURLCache{},
	}, &finalMessages
}

type NewConfigArgs struct {
	ConfigFile   Option[configdomain.PartialConfig]
	DryRun       bool
	GlobalConfig configdomain.PartialConfig
	LocalConfig  configdomain.PartialConfig
	Runner       gitconfig.Runner
}
