package config

import (
	"fmt"
	"slices"

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
	File       Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun     configdomain.DryRun                // whether to only print the Git commands but not execute them
	Env        configdomain.PartialConfig         // configuration data taken from environment variables
	Git        configdomain.PartialConfig         // configuration data taken from Git metadata, in particular the unscoped Git metadata
	GitVersion git.Version                        // version of the installed Git executable
}

// removes the given branch from the lineage, and updates its children
func (self *NormalConfig) CleanupBranchFromLineage(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) {
	parent, hasParent := self.Git.Lineage.Parent(branch).Get()
	children := self.Lineage.Children(branch)
	for _, child := range children {
		if hasParent {
			self.Lineage = self.Lineage.Set(child, parent)
			_ = gitconfig.SetParent(runner, child, parent)
		} else {
			self.Lineage = self.Lineage.RemoveBranch(child)
			_ = gitconfig.RemoveParent(runner, parent)
		}
	}
	self.Lineage = self.Lineage.RemoveBranch(branch)
	_ = gitconfig.RemoveParent(runner, branch)
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
	_ = gitconfig.RemoveBranchTypeOverride(runner, branch)
	return nil
}

func (self *NormalConfig) RemoveDevRemote(runner subshelldomain.Runner) {
	if self.Git.DevRemote.IsSome() {
		_ = gitconfig.RemoveDevRemote(runner)
	}
}

func (self *NormalConfig) RemoveFeatureRegex(runner subshelldomain.Runner) {
	if self.Git.FeatureRegex.IsSome() {
		_ = gitconfig.RemoveFeatureRegex(runner)
	}
}

func (self *NormalConfig) RemoveNewBranchType(runner subshelldomain.Runner) {
	if self.Git.NewBranchType.IsSome() {
		_ = gitconfig.RemoveNewBranchType(runner)
	}
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *NormalConfig) RemoveParent(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) {
	self.Git.Lineage = self.Git.Lineage.RemoveBranch(branch)
	_ = gitconfig.RemoveParent(runner, branch)
}

func (self *NormalConfig) RemovePerennialAncestors(runner subshelldomain.Runner, finalMessages stringslice.Collector) {
	for _, perennialBranch := range self.PerennialBranches {
		if self.Lineage.Parent(perennialBranch).IsSome() {
			_ = gitconfig.RemoveParent(runner, perennialBranch)
			self.Lineage = self.Lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
}

func (self *NormalConfig) RemovePerennialBranches(runner subshelldomain.Runner) {
	if len(self.Git.PerennialBranches) > 0 {
		_ = gitconfig.RemovePerennialBranches(runner)
	}
}

func (self *NormalConfig) RemovePerennialRegex(runner subshelldomain.Runner) {
	if self.Git.PerennialRegex.IsSome() {
		_ = gitconfig.RemovePerennialRegex(runner)
	}
}

func (self *NormalConfig) RemovePushHook(runner subshelldomain.Runner) {
	if self.Git.PushHook.IsSome() {
		_ = gitconfig.RemovePushHook(runner)
	}
}

func (self *NormalConfig) RemoveShareNewBranches(runner subshelldomain.Runner) {
	if self.Git.ShareNewBranches.IsSome() {
		_ = gitconfig.RemoveShareNewBranches(runner)
	}
}

func (self *NormalConfig) RemoveShipDeleteTrackingBranch(runner subshelldomain.Runner) {
	if self.Git.ShipDeleteTrackingBranch.IsSome() {
		_ = gitconfig.RemoveShipDeleteTrackingBranch(runner)
	}
}

func (self *NormalConfig) RemoveShipStrategy(runner subshelldomain.Runner) {
	if self.Git.ShipStrategy.IsSome() {
		_ = gitconfig.RemoveShipStrategy(runner)
	}
}

func (self *NormalConfig) RemoveSyncFeatureStrategy(runner subshelldomain.Runner) {
	if self.Git.SyncFeatureStrategy.IsSome() {
		_ = gitconfig.RemoveSyncFeatureStrategy(runner)
	}
}

func (self *NormalConfig) RemoveSyncPerennialStrategy(runner subshelldomain.Runner) {
	if self.Git.SyncPerennialStrategy.IsSome() {
		_ = gitconfig.RemoveSyncPerennialStrategy(runner)
	}
}

func (self *NormalConfig) RemoveSyncPrototypeStrategy(runner subshelldomain.Runner) {
	if self.Git.SyncPrototypeStrategy.IsSome() {
		_ = gitconfig.RemoveSyncPrototypeStrategy(runner)
	}
}

func (self *NormalConfig) RemoveSyncTags(runner subshelldomain.Runner) {
	if self.Git.SyncTags.IsSome() {
		_ = gitconfig.RemoveSyncTags(runner)
	}
}

func (self *NormalConfig) RemoveSyncUpstream(runner subshelldomain.Runner) {
	if self.Git.SyncUpstream.IsSome() {
		_ = gitconfig.RemoveSyncUpstream(runner)
	}
}

// SetBranchTypeOverride registers the given branch names as contribution branches.
// The branches must exist.
func (self *NormalConfig) SetBranchTypeOverride(runner subshelldomain.Runner, branchType configdomain.BranchType, branches ...gitdomain.LocalBranchName) error {
	for _, branch := range branches {
		self.BranchTypeOverrides[branch] = branchType
		if err := gitconfig.SetBranchTypeOverride(runner, branch, branchType); err != nil {
			return err
		}
	}
	return nil
}

// SetDevRemote updates the locally configured development remote.
func (self *NormalConfig) SetDevRemote(runner subshelldomain.Runner, remote gitdomain.Remote) error {
	self.DevRemote = remote
	existing, has := self.Git.DevRemote.Get()
	if has && existing == remote {
		return nil
	}
	return gitconfig.SetDevRemote(runner, remote)
}

// SetFeatureRegex updates the locally configured feature regex.
func (self *NormalConfig) SetFeatureRegex(runner subshelldomain.Runner, value configdomain.FeatureRegex) error {
	self.FeatureRegex = Some(value)
	existing, has := self.Git.FeatureRegex.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetFeatureRegex(runner, value)
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *NormalConfig) SetNewBranchType(runner subshelldomain.Runner, value configdomain.BranchType) error {
	self.NewBranchType = Some(value)
	existing, has := self.Git.NewBranchType.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetNewBranchType(runner, value)
}

// SetOffline updates whether Git Town is in offline mode.
func (self *NormalConfig) SetOffline(runner subshelldomain.Runner, value configdomain.Offline) error {
	self.Offline = value
	existing, has := self.Git.Offline.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetOffline(runner, value)
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *NormalConfig) SetParent(runner subshelldomain.Runner, branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage = self.Lineage.Set(branch, parentBranch)
	return gitconfig.SetParent(runner, branch, parentBranch)
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *NormalConfig) SetPerennialBranches(runner subshelldomain.Runner, branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	if slices.Compare(self.Git.PerennialBranches, branches) == 0 {
		return nil
	}
	return gitconfig.SetPerennialBranches(runner, branches)
}

// SetPerennialRegex updates the locally configured perennial regex.
func (self *NormalConfig) SetPerennialRegex(runner subshelldomain.Runner, value configdomain.PerennialRegex) error {
	self.PerennialRegex = Some(value)
	existing, has := self.Git.PerennialRegex.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetPerennialRegex(runner, value)
}

// SetPushHook updates the locally configured push-hook strategy.
func (self *NormalConfig) SetPushHook(runner subshelldomain.Runner, value configdomain.PushHook) error {
	self.PushHook = value
	existing, has := self.Git.PushHook.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetPushHook(runner, value)
}

// SetShareNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *NormalConfig) SetShareNewBranches(runner subshelldomain.Runner, value configdomain.ShareNewBranches) error {
	self.ShareNewBranches = value
	existing, has := self.Git.ShareNewBranches.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetShareNewBranches(runner, value)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *NormalConfig) SetShipDeleteTrackingBranch(runner subshelldomain.Runner, value configdomain.ShipDeleteTrackingBranch) error {
	self.ShipDeleteTrackingBranch = value
	existing, has := self.Git.ShipDeleteTrackingBranch.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetShipDeleteTrackingBranch(runner, value)
}

func (self *NormalConfig) SetShipStrategy(runner subshelldomain.Runner, value configdomain.ShipStrategy) error {
	self.ShipStrategy = value
	existing, has := self.Git.ShipStrategy.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetShipStrategy(runner, value)
}

func (self *NormalConfig) SetSyncFeatureStrategy(runner subshelldomain.Runner, value configdomain.SyncFeatureStrategy) error {
	self.SyncFeatureStrategy = value
	existing, has := self.Git.SyncFeatureStrategy.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetSyncFeatureStrategy(runner, value)
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPerennialStrategy(runner subshelldomain.Runner, strategy configdomain.SyncPerennialStrategy) error {
	self.SyncPerennialStrategy = strategy
	existing, has := self.Git.SyncPerennialStrategy.Get()
	if has && existing == strategy {
		return nil
	}
	return gitconfig.SetSyncPerennialStrategy(runner, strategy)
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPrototypeStrategy(runner subshelldomain.Runner, strategy configdomain.SyncPrototypeStrategy) error {
	self.SyncPrototypeStrategy = strategy
	existing, has := self.Git.SyncPrototypeStrategy.Get()
	if has && existing == strategy {
		return nil
	}
	return gitconfig.SetSyncPrototypeStrategy(runner, strategy)
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncTags(runner subshelldomain.Runner, value configdomain.SyncTags) error {
	self.SyncTags = value
	existing, has := self.Git.SyncTags.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetSyncTags(runner, value)
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *NormalConfig) SetSyncUpstream(runner subshelldomain.Runner, value configdomain.SyncUpstream) error {
	self.SyncUpstream = value
	existing, has := self.Git.SyncUpstream.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetSyncUpstream(runner, value)
}

// SetUnknownBranchType updates the locally configured unknown branch type.
func (self *NormalConfig) SetUnknownBranchType(runner subshelldomain.Runner, value configdomain.BranchType) error {
	self.UnknownBranchType = value
	existing, has := self.Git.UnknownBranchType.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetUnknownBranchType(runner, value)
}
