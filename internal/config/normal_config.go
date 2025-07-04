package config

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type NormalConfig struct {
	configdomain.NormalConfigData
	ConfigFile     Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun         configdomain.DryRun                // whether to only print the Git commands but not execute them
	EnvConfig      configdomain.PartialConfig         // content of the Git Town related environment variables
	GitConfig      configdomain.PartialConfig         // content of the unscoped Git configuration
	GitPersistence gitconfig.Persistence              // access to Git configuration settings
	GitVersion     git.Version                        // version of the installed Git executable
}

// removes the given branch from the lineage, and updates its children
func (self *NormalConfig) CleanupBranchFromLineage(branch gitdomain.LocalBranchName) {
	parent, hasParent := self.GitConfig.Lineage.Parent(branch).Get()
	children := self.Lineage.Children(branch)
	for _, child := range children {
		if hasParent {
			self.Lineage = self.Lineage.Set(child, parent)
			_ = self.GitPersistence.SetParent(child, parent)
		} else {
			self.Lineage = self.Lineage.RemoveBranch(child)
			self.GitPersistence.RemoveParent(parent)
		}
	}
	self.Lineage = self.Lineage.RemoveBranch(branch)
	_ = self.GitPersistence.RemoveParent(branch)
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
	urlStr, hasURLStr := self.RemoteURLString(remote).Get()
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
func (self *NormalConfig) RemoteURLString(remote gitdomain.Remote) Option[string] {
	remoteOverride := envconfig.RemoteURLOverride()
	if remoteOverride.IsSome() {
		return remoteOverride
	}
	return self.GitPersistence.RemoteURL(remote)
}

func (self *NormalConfig) RemoveBranchTypeOverride(branch gitdomain.LocalBranchName) error {
	delete(self.BranchTypeOverrides, branch)
	_ = self.GitPersistence.RemoveBranchTypeOverride(branch)
	return nil
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *NormalConfig) RemoveParent(branch gitdomain.LocalBranchName) {
	self.GitConfig.Lineage = self.GitConfig.Lineage.RemoveBranch(branch)
	_ = self.GitPersistence.RemoveParent(branch)
}

func (self *NormalConfig) RemovePerennialAncestors(finalMessages stringslice.Collector) {
	for _, perennialBranch := range self.PerennialBranches {
		if self.Lineage.Parent(perennialBranch).IsSome() {
			_ = self.GitPersistence.RemoveParent(perennialBranch)
			self.Lineage = self.Lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
}

// SetBranchTypeOverride registers the given branch names as contribution branches.
// The branches must exist.
func (self *NormalConfig) SetBranchTypeOverride(branchType configdomain.BranchType, branches ...gitdomain.LocalBranchName) error {
	for _, branch := range branches {
		self.BranchTypeOverrides[branch] = branchType
		if err := self.GitPersistence.SetBranchTypeOverride(branch, branchType); err != nil {
			return err
		}
	}
	return nil
}

// SetDevRemote updates the locally configured development remote.
func (self *NormalConfig) SetDevRemote(value gitdomain.Remote) error {
	self.DevRemote = value
	return self.GitPersistence.SetDevRemote(value)
}

// SetFeatureRegexLocally updates the locally configured feature regex.
func (self *NormalConfig) SetFeatureRegexLocally(value configdomain.FeatureRegex) error {
	self.FeatureRegex = Some(value)
	return self.GitPersistence.SetFeatureRegex(value)
}

// SetContributionBranches marks the given branches as contribution branches.
func (self *NormalConfig) SetNewBranchType(value configdomain.BranchType) error {
	self.NewBranchType = Some(value)
	return self.GitPersistence.SetNewBranchType(value)
}

// SetOffline updates whether Git Town is in offline mode.
func (self *NormalConfig) SetOffline(value configdomain.Offline) error {
	self.Offline = value
	return self.GitPersistence.SetOffline(value)
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *NormalConfig) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage = self.Lineage.Set(branch, parentBranch)
	return self.GitPersistence.SetParent(branch, parentBranch)
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *NormalConfig) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	return self.GitPersistence.SetPerennialBranches(branches)
}

// SetPerennialRegexLocally updates the locally configured perennial regex.
func (self *NormalConfig) SetPerennialRegexLocally(value configdomain.PerennialRegex) error {
	self.PerennialRegex = Some(value)
	return self.GitPersistence.SetPerennialRegex(value)
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *NormalConfig) SetPushHookLocally(value configdomain.PushHook) error {
	self.PushHook = value
	return self.GitPersistence.SetPushHook(value)
}

// SetShareNewBranches updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *NormalConfig) SetShareNewBranches(value configdomain.ShareNewBranches, scope configdomain.ConfigScope) error {
	self.ShareNewBranches = value
	return self.GitPersistence.SetShareNewBranches(scope, value)
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *NormalConfig) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, scope configdomain.ConfigScope) error {
	self.ShipDeleteTrackingBranch = value
	return self.GitPersistence.SetShipDeleteTrackingBranch(scope, value)
}

func (self *NormalConfig) SetShipStrategy(value configdomain.ShipStrategy, scope configdomain.ConfigScope) error {
	self.ShipStrategy = value
	return self.GitPersistence.SetShipStrategy(scope, value)
}

func (self *NormalConfig) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.SyncFeatureStrategy = value
	return self.GitIO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.SyncPerennialStrategy = strategy
	return self.GitIO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncPrototypeStrategy(strategy configdomain.SyncPrototypeStrategy) error {
	self.SyncPrototypeStrategy = strategy
	return self.GitIO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncPrototypeStrategy, strategy.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *NormalConfig) SetSyncTags(value configdomain.SyncTags) error {
	self.SyncTags = value
	return self.GitIO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeySyncTags, value.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *NormalConfig) SetSyncUpstream(value configdomain.SyncUpstream, scope configdomain.ConfigScope) error {
	self.SyncUpstream = value
	return self.GitIO.SetConfigValue(scope, configdomain.KeySyncUpstream, strconv.FormatBool(value.IsTrue()))
}

// SetUnknownBranchTypeLocally updates the locally configured unknown branch type.
func (self *NormalConfig) SetUnknownBranchTypeLocally(value configdomain.BranchType) error {
	self.UnknownBranchType = value
	return self.GitIO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyUnknownBranchType, value.String())
}
