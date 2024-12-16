package config

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/config/confighelpers"
	"github.com/git-town/git-town/v16/internal/config/envconfig"
	"github.com/git-town/git-town/v16/internal/config/gitconfig"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/git/giturl"
	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type NormalConfig struct {
	configdomain.NormalConfigData
	ConfigFile      Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun          configdomain.DryRun                // whether to only print the Git commands but not execute them
	GitConfigAccess gitconfig.Access                   // access to the Git configuration settings
	GitVersion      git.Version                        // version of the installed Git executable
	GlobalGitConfig configdomain.PartialConfig         // content of the global Git configuration
	LocalGitConfig  configdomain.PartialConfig         // content of the local Git configuration
}

// AddToContributionBranches registers the given branch names as contribution branches.
// The branches must exist.
func (self *NormalConfig) AddToContributionBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetContributionBranches(append(self.ContributionBranches, branches...))
}

// AddToObservedBranches registers the given branch names as observed branches.
// The branches must exist.
func (self *NormalConfig) AddToObservedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetObservedBranches(append(self.ObservedBranches, branches...))
}

// AddToParkedBranches registers the given branch names as parked branches.
// The branches must exist.
func (self *NormalConfig) AddToParkedBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetParkedBranches(append(self.ParkedBranches, branches...))
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *NormalConfig) AddToPerennialBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetPerennialBranches(append(self.PerennialBranches, branches...))
}

// AddToPrototypeBranches registers the given branch names as prototype branches.
// The branches must exist.
func (self *NormalConfig) AddToPrototypeBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetPrototypeBranches(append(self.PrototypeBranches, branches...))
}

// removes the given branch from the lineage, and updates its children
func (self *NormalConfig) CleanupBranchFromLineage(branch gitdomain.LocalBranchName) {
	parent, hasParent := self.LocalGitConfig.Lineage.Parent(branch).Get()
	children := self.Lineage.Children(branch)
	for _, child := range children {
		if hasParent {
			self.Lineage = self.Lineage.Set(child, parent)
			_ = self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
		} else {
			self.Lineage = self.Lineage.RemoveBranch(child)
			_ = self.GitConfigAccess.RemoveConfigValue(configdomain.ConfigScopeLocal, configdomain.NewParentKey(parent))
		}
	}
	self.Lineage = self.Lineage.RemoveBranch(branch)
	_ = self.GitConfigAccess.RemoveConfigValue(configdomain.ConfigScopeLocal, configdomain.NewParentKey(branch))
}

// DevURL provides the URL for the development remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) DevURL() Option[giturl.Parts] {
	return self.RemoteURL(self.DevRemote)
}

// RemoteURL provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) RemoteURL(remote gitdomain.Remote) Option[giturl.Parts] {
	text, hasText := self.RemoteURLString(remote).Get()
	if !hasText {
		return None[giturl.Parts]()
	}
	return confighelpers.DetermineRemoteURL(text, self.HostingOriginHostname)
}

// RemoteURLString provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *NormalConfig) RemoteURLString(remote gitdomain.Remote) Option[string] {
	remoteOverride := envconfig.RemoteURLOverride()
	if remoteOverride.IsSome() {
		return remoteOverride
	}
	return self.GitConfigAccess.RemoteURL(remote)
}

func (self *NormalConfig) RemoveCreatePrototypeBranches() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyDeprecatedCreatePrototypeBranches)
}

func (self *NormalConfig) RemoveDevRemote() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyDevRemote)
}

func (self *NormalConfig) RemoveFeatureRegex() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyFeatureRegex)
}

// RemoveFromContributionBranches removes the given branch as a perennial branch.
func (self *NormalConfig) RemoveFromContributionBranches(branch gitdomain.LocalBranchName) error {
	self.ContributionBranches = slice.Remove(self.ContributionBranches, branch)
	return self.SetContributionBranches(self.ContributionBranches)
}

// RemoveFromObservedBranches removes the given branch as a perennial branch.
func (self *NormalConfig) RemoveFromObservedBranches(branch gitdomain.LocalBranchName) error {
	self.ObservedBranches = slice.Remove(self.ObservedBranches, branch)
	return self.SetObservedBranches(self.ObservedBranches)
}

// RemoveFromParkedBranches removes the given branch as a perennial branch.
func (self *NormalConfig) RemoveFromParkedBranches(branch gitdomain.LocalBranchName) error {
	self.ParkedBranches = slice.Remove(self.ParkedBranches, branch)
	return self.SetParkedBranches(self.ParkedBranches)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *NormalConfig) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	self.PerennialBranches = slice.Remove(self.PerennialBranches, branch)
	return self.SetPerennialBranches(self.PerennialBranches)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *NormalConfig) RemoveFromPrototypeBranches(branch gitdomain.LocalBranchName) error {
	self.PrototypeBranches = slice.Remove(self.PrototypeBranches, branch)
	return self.SetPrototypeBranches(self.PrototypeBranches)
}

func (self *NormalConfig) RemoveNewBranchType() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyNewBranchType)
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *NormalConfig) RemoveParent(branch gitdomain.LocalBranchName) {
	self.LocalGitConfig.Lineage = self.LocalGitConfig.Lineage.RemoveBranch(branch)
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.NewParentKey(branch))
}

func (self *NormalConfig) RemovePerennialAncestors(finalMessages stringslice.Collector) {
	for _, perennialBranch := range self.PerennialBranches {
		if self.Lineage.Parent(perennialBranch).IsSome() {
			_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.NewParentKey(perennialBranch))
			self.Lineage = self.Lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
}

func (self *NormalConfig) RemovePerennialBranches() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyPerennialBranches)
}

func (self *NormalConfig) RemovePerennialRegex() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyPerennialRegex)
}

func (self *NormalConfig) RemovePushHook() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyPushHook)
}

func (self *NormalConfig) RemovePushNewBranches() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyPushNewBranches)
}

func (self *NormalConfig) RemoveShipDeleteTrackingBranch() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyShipDeleteTrackingBranch)
}

func (self *NormalConfig) RemoveShipStrategy() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyShipStrategy)
}

func (self *NormalConfig) RemoveSyncFeatureStrategy() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeySyncFeatureStrategy)
}

func (self *NormalConfig) RemoveSyncPerennialStrategy() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeySyncPerennialStrategy)
}

func (self *NormalConfig) RemoveSyncPrototypeStrategy() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeySyncPrototypeStrategy)
}

func (self *NormalConfig) RemoveSyncTags() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeySyncTags)
}

func (self *NormalConfig) RemoveSyncUpstream() {
	_ = self.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeySyncUpstream)
}

// SetObservedBranches marks the given branches as observed branches.
func (self *NormalConfig) SetContributionBranches(branches gitdomain.LocalBranchNames) error {
	self.ContributionBranches = branches
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyContributionBranches, branches.Join(" "))
}

// SetDefaultBranchTypeLocally updates the locally configured default branch type.
func (self *NormalConfig) SetDefaultBranchTypeLocally(value configdomain.BranchType) error {
	self.DefaultBranchType = value
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyDefaultBranchType, value.String())
}

// SetDefaultBranchTypeLocally updates the locally configured default branch type.
func (self *NormalConfig) SetDevRemote(value gitdomain.Remote) error {
	self.DevRemote = value
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyDevRemote, value.String())
}

// SetFeatureRegexLocally updates the locally configured feature regex.
func (self *NormalConfig) SetFeatureRegexLocally(value configdomain.FeatureRegex) error {
	self.FeatureRegex = Some(value)
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex, value.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *NormalConfig) SetNewBranchType(value configdomain.BranchType) error {
	self.NewBranchType = value
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType, value.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *NormalConfig) SetObservedBranches(branches gitdomain.LocalBranchNames) error {
	self.ObservedBranches = branches
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyObservedBranches, branches.Join(" "))
}

// SetOffline updates whether Git Town is in offline mode.
func (self *NormalConfig) SetOffline(value configdomain.Offline) error {
	self.Offline = value
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeGlobal, configdomain.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *NormalConfig) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage = self.Lineage.Set(branch, parentBranch)
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.NewParentKey(branch), parentBranch.String())
}

// SetObservedBranches marks the given branches as perennial branches.
func (self *NormalConfig) SetParkedBranches(branches gitdomain.LocalBranchNames) error {
	self.ParkedBranches = branches
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyParkedBranches, branches.Join(" "))
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *NormalConfig) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyPerennialBranches, branches.Join(" "))
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *NormalConfig) SetPerennialRegexLocally(value configdomain.PerennialRegex) error {
	self.PerennialRegex = Some(value)
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyPerennialRegex, value.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *NormalConfig) SetPrototypeBranches(branches gitdomain.LocalBranchNames) error {
	self.PrototypeBranches = branches
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyPrototypeBranches, branches.Join(" "))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *NormalConfig) SetPushHookLocally(value configdomain.PushHook) error {
	self.PushHook = value
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetPushNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *NormalConfig) SetPushNewBranches(value configdomain.PushNewBranches, scope configdomain.ConfigScope) error {
	setting := strconv.FormatBool(bool(value))
	self.PushNewBranches = value
	return self.GitConfigAccess.SetConfigValue(scope, configdomain.KeyPushNewBranches, setting)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *NormalConfig) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, scope configdomain.ConfigScope) error {
	self.ShipDeleteTrackingBranch = value
	return self.GitConfigAccess.SetConfigValue(scope, configdomain.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.IsTrue()))
}

func (self *NormalConfig) SetShipStrategy(value configdomain.ShipStrategy, scope configdomain.ConfigScope) error {
	self.ShipStrategy = value
	return self.GitConfigAccess.SetConfigValue(scope, configdomain.KeyShipStrategy, value.String())
}

func (self *NormalConfig) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.SyncFeatureStrategy = value
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.SyncPerennialStrategy = strategy
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPrototypeStrategy(strategy configdomain.SyncPrototypeStrategy) error {
	self.SyncPrototypeStrategy = strategy
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncPrototypeStrategy, strategy.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncTags(value configdomain.SyncTags) error {
	self.SyncTags = value
	return self.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncTags, value.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *NormalConfig) SetSyncUpstream(value configdomain.SyncUpstream, scope configdomain.ConfigScope) error {
	self.SyncUpstream = value
	return self.GitConfigAccess.SetConfigValue(scope, configdomain.KeySyncUpstream, strconv.FormatBool(value.IsTrue()))
}
