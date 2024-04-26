// Package config provides functionality to read and write the Git Town configuration.
// Git Town configuration can exist in a number of locations: in local or global Git metadata or in a configuration file.
// Subspackages implement access to specific configuration locations.
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
	"github.com/git-town/git-town/v14/src/messages"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type Config struct {
	ConfigFile      *configdomain.PartialConfig // content of git-town.toml, nil = no config file exists
	DryRun          bool
	FullConfig      configdomain.FullConfig    // the merged configuration data
	GitConfig       gitconfig.Access           // access to the Git configuration settings
	GlobalGitConfig configdomain.PartialConfig // content of the global Git configuration
	LocalGitConfig  configdomain.PartialConfig // content of the local Git configuration
	originURLCache  configdomain.OriginURLCache
}

// AddToContributionBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToContributionBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetContributionBranches(append(self.FullConfig.ContributionBranches, branches...))
}

// AddToObservedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToObservedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetObservedBranches(append(self.FullConfig.ObservedBranches, branches...))
}

// AddToParkedBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToParkedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetParkedBranches(append(self.FullConfig.ParkedBranches, branches...))
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToPerennialBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetPerennialBranches(append(self.FullConfig.PerennialBranches, branches...))
}

// Author provides the locally Git configured user.
func (self *Config) Author() gitdomain.Author {
	email := self.FullConfig.GitUserEmail
	name := self.FullConfig.GitUserName
	return gitdomain.Author(fmt.Sprintf("%s <%s>", name, email))
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *Config) OriginURL() *giturl.Parts {
	text := self.OriginURLString()
	if text == "" {
		return nil
	}
	return confighelpers.DetermineOriginURL(text, self.FullConfig.HostingOriginHostname, self.originURLCache)
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
	self.FullConfig = configdomain.NewFullConfig(self.ConfigFile, self.GlobalGitConfig, self.LocalGitConfig)
}

// RemoveFromContributionBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromContributionBranches(branch gitdomain.LocalBranchName) error {
	self.FullConfig.ContributionBranches = slice.Remove(self.FullConfig.ContributionBranches, branch)
	return self.SetContributionBranches(self.FullConfig.ContributionBranches)
}

// RemoveFromObservedBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromObservedBranches(branch gitdomain.LocalBranchName) error {
	self.FullConfig.ObservedBranches = slice.Remove(self.FullConfig.ObservedBranches, branch)
	return self.SetObservedBranches(self.FullConfig.ObservedBranches)
}

// RemoveFromParkedBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromParkedBranches(branch gitdomain.LocalBranchName) error {
	self.FullConfig.ParkedBranches = slice.Remove(self.FullConfig.ParkedBranches, branch)
	return self.SetParkedBranches(self.FullConfig.ParkedBranches)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	self.FullConfig.PerennialBranches = slice.Remove(self.FullConfig.PerennialBranches, branch)
	return self.SetPerennialBranches(self.FullConfig.PerennialBranches)
}

func (self *Config) RemoveMainBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.KeyMainBranch)
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (self *Config) RemoveOutdatedConfiguration(localBranches gitdomain.LocalBranchNames) error {
	for child, parent := range self.FullConfig.Lineage {
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
	self.FullConfig.ContributionBranches = branches
	self.LocalGitConfig.ContributionBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyContributionBranches, branches.Join(" "))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *Config) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.FullConfig.MainBranch = branch
	self.LocalGitConfig.MainBranch = Some(branch)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, branch.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *Config) SetObservedBranches(branches gitdomain.LocalBranchNames) error {
	self.FullConfig.ObservedBranches = branches
	self.LocalGitConfig.ObservedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyObservedBranches, branches.Join(" "))
}

// SetOffline updates whether Git Town is in offline mode.
func (self *Config) SetOffline(value configdomain.Offline) error {
	self.FullConfig.Offline = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyOffline, value.String())
}

// SetOriginHostname marks the given branch as the main branch
// in the Git Town configuration.
func (self *Config) SetOriginHostname(hostName configdomain.HostingOriginHostname) error {
	self.FullConfig.HostingOriginHostname = Some(hostName)
	self.LocalGitConfig.HostingOriginHostname = Some(hostName)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyHostingOriginHostname, hostName.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *Config) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.FullConfig.Lineage[branch] = parentBranch
	return self.GitConfig.SetLocalConfigValue(gitconfig.NewParentKey(branch), parentBranch.String())
}

// SetObservedBranches marks the given branches as perennial branches.
func (self *Config) SetParkedBranches(branches gitdomain.LocalBranchNames) error {
	self.FullConfig.ParkedBranches = branches
	self.LocalGitConfig.ParkedBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyParkedBranches, branches.Join(" "))
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *Config) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.FullConfig.PerennialBranches = branches
	self.LocalGitConfig.PerennialBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, branches.Join(" "))
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *Config) SetPerennialRegexLocally(value configdomain.PerennialRegex) error {
	self.LocalGitConfig.PerennialRegex = Some(value)
	self.FullConfig.PerennialRegex = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialRegex, value.String())
}

// SetPushHook updates the configured push-hook strategy.
func (self *Config) SetPushHookGlobally(value configdomain.PushHook) error {
	self.GlobalGitConfig.PushHook = Some(value)
	self.FullConfig.PushHook = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(value.Bool()))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *Config) SetPushHookLocally(value configdomain.PushHook) error {
	self.LocalGitConfig.PushHook = Some(value)
	self.FullConfig.PushHook = value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetPushNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *Config) SetPushNewBranches(value configdomain.PushNewBranches, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.FullConfig.PushNewBranches = value
	if global {
		self.GlobalGitConfig.PushNewBranches = Some(value)
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushNewBranches, setting)
	}
	self.LocalGitConfig.PushNewBranches = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushNewBranches, setting)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *Config) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, global bool) error {
	self.FullConfig.ShipDeleteTrackingBranch = value
	if global {
		self.GlobalGitConfig.ShipDeleteTrackingBranch = Some(value)
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
	}
	self.LocalGitConfig.ShipDeleteTrackingBranch = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *Config) SetSyncBeforeShip(value configdomain.SyncBeforeShip, global bool) error {
	self.FullConfig.SyncBeforeShip = value
	if global {
		self.GlobalGitConfig.SyncBeforeShip = Some(value)
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
	}
	self.LocalGitConfig.SyncBeforeShip = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
}

func (self *Config) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.FullConfig.SyncFeatureStrategy = value
	self.LocalGitConfig.SyncFeatureStrategy = Some(value)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

func (self *Config) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	self.GlobalGitConfig.SyncFeatureStrategy = Some(value)
	self.FullConfig.SyncFeatureStrategy = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *Config) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.FullConfig.SyncPerennialStrategy = strategy
	self.LocalGitConfig.SyncPerennialStrategy = Some(strategy)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *Config) SetSyncUpstream(value configdomain.SyncUpstream, global bool) error {
	self.FullConfig.SyncUpstream = value
	if global {
		self.GlobalGitConfig.SyncUpstream = Some(value)
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}

func NewConfig(args NewConfigArgs) (*Config, *stringslice.Collector, error) {
	config := configdomain.NewFullConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	configAccess := gitconfig.Access{Runner: args.Runner}
	finalMessages := stringslice.Collector{}
	err := cleanupPerennialParentEntries(config.Lineage, config.MainAndPerennials(), configAccess, &finalMessages)
	return &Config{
		ConfigFile:      args.ConfigFile,
		DryRun:          args.DryRun,
		FullConfig:      config,
		GitConfig:       configAccess,
		GlobalGitConfig: args.GlobalConfig,
		LocalGitConfig:  args.LocalConfig,
		originURLCache:  configdomain.OriginURLCache{},
	}, &finalMessages, err
}

type NewConfigArgs struct {
	ConfigFile   *configdomain.PartialConfig
	DryRun       bool
	GlobalConfig configdomain.PartialConfig
	LocalConfig  configdomain.PartialConfig
	Runner       gitconfig.Runner
}

// cleanupPerennialParentEntries removes outdated entries from the configuration.
func cleanupPerennialParentEntries(lineage configdomain.Lineage, perennialBranches gitdomain.LocalBranchNames, access gitconfig.Access, finalMessages *stringslice.Collector) error {
	for _, perennialBranch := range perennialBranches {
		if lineage.Parent(perennialBranch).IsSome() {
			if err := access.RemoveLocalConfigValue(gitconfig.NewParentKey(perennialBranch)); err != nil {
				return err
			}
			lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
	return nil
}
