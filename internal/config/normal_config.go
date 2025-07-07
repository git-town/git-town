package config

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type NormalConfig struct {
	configdomain.NormalConfigData
	ConfigFile Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun     configdomain.DryRun                // whether to only print the Git commands but not execute them
	EnvConfig  configdomain.PartialConfig         // content of the Git Town related environment variables
	GitConfig  configdomain.PartialConfig         // content of the unscoped Git configuration
	GitVersion git.Version                        // version of the installed Git executable
}

// removes the given branch from the lineage, and updates its children
func (self *NormalConfig) CleanupBranchFromLineage(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) {
	parent, hasParent := self.GitConfig.Lineage.Parent(branch).Get()
	children := self.Lineage.Children(branch)
	for _, child := range children {
		if hasParent {
			self.Lineage = self.Lineage.Set(child, parent)
			_ = gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
		} else {
			self.Lineage = self.Lineage.RemoveBranch(child)
			_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(parent))
		}
	}
	self.Lineage = self.Lineage.RemoveBranch(branch)
	_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(branch))
}

// DevURL provides the URL for the development remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) DevURL(querier subshelldomain.Querier) Option[giturl.Parts] {
	return self.RemoteURL(querier, self.DevRemote)
}

// RemoteURL provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) RemoteURL(querier subshelldomain.Querier, remote gitdomain.Remote) Option[giturl.Parts] {
	urlStr, hasURLStr := self.RemoteURLString(querier, remote).Get()
	if !hasURLStr {
		return None[giturl.Parts]()
	}
	url, hasURL := giturl.Parse(urlStr).Get()
	if !hasURL {
		return None[giturl.Parts]()
	}
	if hostnameOverride, hasHostNameOverride := self.HostingOriginHostname.Get(); hasHostNameOverride {
		url.Host = hostnameOverride.String()
	}
	return Some(url)
}

// RemoteURLString provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *NormalConfig) RemoteURLString(querier subshelldomain.Querier, remote gitdomain.Remote) Option[string] {
	remoteOverride := envconfig.RemoteURLOverride()
	if remoteOverride.IsSome() {
		return remoteOverride
	}
	return gitconfig.RemoteURL(querier, remote)
}

func (self *NormalConfig) RemoveBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	delete(self.BranchTypeOverrides, branch)
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
	return nil
}

func (self *NormalConfig) RemoveDevRemote(runner subshelldomain.Runner) {
	if self.GitConfig.DevRemote.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDevRemote)
	}
}

func (self *NormalConfig) RemoveFeatureRegex(runner subshelldomain.Runner) {
	if self.GitConfig.FeatureRegex.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex)
	}
}

func (self *NormalConfig) RemoveNewBranchType(runner subshelldomain.Runner) {
	if self.GitConfig.NewBranchType.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType)
	}
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *NormalConfig) RemoveParent(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) {
	self.GitConfig.Lineage = self.GitConfig.Lineage.RemoveBranch(branch)
	_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(branch))
}

func (self *NormalConfig) RemovePerennialAncestors(runner subshelldomain.Runner, finalMessages stringslice.Collector) {
	for _, perennialBranch := range self.PerennialBranches {
		if self.Lineage.Parent(perennialBranch).IsSome() {
			_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(perennialBranch))
			self.Lineage = self.Lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
}

func (self *NormalConfig) RemovePerennialBranches(runner subshelldomain.Runner) {
	if len(self.GitConfig.PerennialBranches) > 0 {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialBranches)
	}
}

func (self *NormalConfig) RemovePerennialRegex(runner subshelldomain.Runner) {
	if self.GitConfig.PerennialRegex.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialRegex)
	}
}

func (self *NormalConfig) RemovePushHook(runner subshelldomain.Runner) {
	if self.GitConfig.PushHook.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPushHook)
	}
}

func (self *NormalConfig) RemoveShareNewBranches(runner subshelldomain.Runner) {
	if self.GitConfig.ShareNewBranches.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShareNewBranches)
	}
}

func (self *NormalConfig) RemoveShipDeleteTrackingBranch(runner subshelldomain.Runner) {
	if self.GitConfig.ShipDeleteTrackingBranch.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipDeleteTrackingBranch)
	}
}

func (self *NormalConfig) RemoveShipStrategy(runner subshelldomain.Runner) {
	if self.GitConfig.ShipStrategy.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyShipStrategy)
	}
}

func (self *NormalConfig) RemoveSyncFeatureStrategy(runner subshelldomain.Runner) {
	if self.GitConfig.SyncFeatureStrategy.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncFeatureStrategy)
	}
}

func (self *NormalConfig) RemoveSyncPerennialStrategy(runner subshelldomain.Runner) {
	if self.GitConfig.SyncPerennialStrategy.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPerennialStrategy)
	}
}

func (self *NormalConfig) RemoveSyncPrototypeStrategy(runner subshelldomain.Runner) {
	if self.GitConfig.SyncPrototypeStrategy.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPrototypeStrategy)
	}
}

func (self *NormalConfig) RemoveSyncTags(runner subshelldomain.Runner) {
	if self.GitConfig.SyncTags.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncTags)
	}
}

func (self *NormalConfig) RemoveSyncUpstream(runner subshelldomain.Runner) {
	if self.GitConfig.SyncUpstream.IsSome() {
		_ = gitconfig.RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncUpstream)
	}
}

// SetBranchTypeOverride registers the given branch names as contribution branches.
// The branches must exist.
func (self *NormalConfig) SetBranchTypeOverride(runner subshelldomain.Runner, branchType configdomain.BranchType, branches ...gitdomain.LocalBranchName) error {
	for _, branch := range branches {
		self.BranchTypeOverrides[branch] = branchType
		if err := gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewBranchTypeOverrideKeyForBranch(branch).Key, branchType.String()); err != nil {
			return err
		}
	}
	return nil
}

// SetDevRemote updates the locally configured development remote.
func (self *NormalConfig) SetDevRemote(runner subshelldomain.Runner, value gitdomain.Remote) error {
	self.DevRemote = value
	existing, has := self.GitConfig.DevRemote.Get()
	if has || existing == value {
		return nil
	}
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDevRemote, value.String())
}

// SetFeatureRegexLocally updates the locally configured feature regex.
func (self *NormalConfig) SetFeatureRegexLocally(runner subshelldomain.Runner, value configdomain.FeatureRegex) error {
	self.FeatureRegex = Some(value)
	existing, has := self.GitConfig.FeatureRegex.Get()
	if has || existing == value {
		return nil
	}
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyFeatureRegex, value.String())
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *NormalConfig) SetNewBranchType(runner subshelldomain.Runner, value configdomain.BranchType) error {
	self.NewBranchType = Some(value)
	existing, has := self.GitConfig.NewBranchType.Get()
	if has || existing == value {
		return nil
	}
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyNewBranchType, value.String())
}

// SetOffline updates whether Git Town is in offline mode.
func (self *NormalConfig) SetOffline(runner subshelldomain.Runner, value configdomain.Offline) error {
	self.Offline = value
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeGlobal, configdomain.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *NormalConfig) SetParent(runner subshelldomain.Runner, branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage = self.Lineage.Set(branch, parentBranch)
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(branch), parentBranch.String())
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *NormalConfig) SetPerennialBranches(runner subshelldomain.Runner, branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	if slices.Compare(self.GitConfig.PerennialBranches, branches) == 0 {
		return nil
	}
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialBranches, branches.Join(" "))
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *NormalConfig) SetPerennialRegexLocally(runner subshelldomain.Runner, value configdomain.PerennialRegex) error {
	self.PerennialRegex = Some(value)
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPerennialRegex, value.String())
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *NormalConfig) SetPushHookLocally(runner subshelldomain.Runner, value configdomain.PushHook) error {
	self.PushHook = value
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetShareNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *NormalConfig) SetShareNewBranches(runner subshelldomain.Runner, value configdomain.ShareNewBranches, scope configdomain.ConfigScope) error {
	self.ShareNewBranches = value
	return gitconfig.SetConfigValue(runner, scope, configdomain.KeyShareNewBranches, value.String())
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *NormalConfig) SetShipDeleteTrackingBranch(runner subshelldomain.Runner, value configdomain.ShipDeleteTrackingBranch, scope configdomain.ConfigScope) error {
	self.ShipDeleteTrackingBranch = value
	return gitconfig.SetConfigValue(runner, scope, configdomain.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.IsTrue()))
}

func (self *NormalConfig) SetShipStrategy(runner subshelldomain.Runner, value configdomain.ShipStrategy, scope configdomain.ConfigScope) error {
	self.ShipStrategy = value
	return gitconfig.SetConfigValue(runner, scope, configdomain.KeyShipStrategy, value.String())
}

func (self *NormalConfig) SetSyncFeatureStrategy(runner subshelldomain.Runner, value configdomain.SyncFeatureStrategy) error {
	self.SyncFeatureStrategy = value
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPerennialStrategy(runner subshelldomain.Runner, strategy configdomain.SyncPerennialStrategy) error {
	self.SyncPerennialStrategy = strategy
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPrototypeStrategy(runner subshelldomain.Runner, strategy configdomain.SyncPrototypeStrategy) error {
	self.SyncPrototypeStrategy = strategy
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncPrototypeStrategy, strategy.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncTags(runner subshelldomain.Runner, value configdomain.SyncTags) error {
	self.SyncTags = value
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeySyncTags, value.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *NormalConfig) SetSyncUpstream(runner subshelldomain.Runner, value configdomain.SyncUpstream, scope configdomain.ConfigScope) error {
	self.SyncUpstream = value
	return gitconfig.SetConfigValue(runner, scope, configdomain.KeySyncUpstream, strconv.FormatBool(value.IsTrue()))
}

// SetUnknownBranchTypeLocally updates the locally configured unknown branch type.
func (self *NormalConfig) SetUnknownBranchTypeLocally(runner subshelldomain.Runner, value configdomain.BranchType) error {
	self.UnknownBranchType = value
	return gitconfig.SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyUnknownBranchType, value.String())
}
