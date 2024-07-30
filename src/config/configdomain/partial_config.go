package configdomain

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                  Aliases
	ContributionBranches     gitdomain.LocalBranchNames
	CreatePrototypeBranches  Option[CreatePrototypeBranches]
	GitHubToken              Option[GitHubToken]
	GitLabToken              Option[GitLabToken]
	GitUserEmail             Option[GitUserEmail]
	GitUserName              Option[GitUserName]
	GiteaToken               Option[GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	HostingPlatform          Option[HostingPlatform]
	Lineage                  Lineage
	MainBranch               Option[gitdomain.LocalBranchName]
	ObservedBranches         gitdomain.LocalBranchNames
	Offline                  Option[Offline]
	ParkedBranches           gitdomain.LocalBranchNames
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PrototypeBranches        gitdomain.LocalBranchNames
	PushHook                 Option[PushHook]
	PushNewBranches          Option[PushNewBranches]
	ShipDeleteTrackingBranch Option[ShipDeleteTrackingBranch]
	SyncFeatureStrategy      Option[SyncFeatureStrategy]
	SyncPerennialStrategy    Option[SyncPerennialStrategy]
	SyncUpstream             Option[SyncUpstream]
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{
		Aliases: Aliases{},
	} //exhaustruct:ignore
}

// Merges the given PartialConfig into this configuration object.
// TODO: refactor to have the same structure as ToUnvalidatedConfig
func (self *PartialConfig) Merge(other PartialConfig) {
	for key, value := range other.Aliases {
		self.Aliases[key] = value
	}
	for _, entry := range other.Lineage.Entries() {
		self.Lineage.Add(entry.Child, entry.Parent)
	}
	self.ContributionBranches = append(self.ContributionBranches, other.ContributionBranches...)
	if other.CreatePrototypeBranches.IsSome() {
		self.CreatePrototypeBranches = other.CreatePrototypeBranches
	}
	if other.HostingOriginHostname.IsSome() {
		self.HostingOriginHostname = other.HostingOriginHostname
	}
	if other.HostingPlatform.IsSome() {
		self.HostingPlatform = other.HostingPlatform
	}
	if other.GiteaToken.IsSome() {
		self.GiteaToken = other.GiteaToken
	}
	if other.GitHubToken.IsSome() {
		self.GitHubToken = other.GitHubToken
	}
	if other.GitLabToken.IsSome() {
		self.GitLabToken = other.GitLabToken
	}
	if other.GitUserEmail.IsSome() {
		self.GitUserEmail = other.GitUserEmail
	}
	if other.GitUserName.IsSome() {
		self.GitUserName = other.GitUserName
	}
	if branch, has := other.MainBranch.Get(); has {
		self.MainBranch = Some(branch)
	}
	if other.PushNewBranches.IsSome() {
		self.PushNewBranches = other.PushNewBranches
	}
	self.ObservedBranches = append(self.ObservedBranches, other.ObservedBranches...)
	if other.Offline.IsSome() {
		self.Offline = other.Offline
	}
	self.ParkedBranches = append(self.ParkedBranches, other.ParkedBranches...)
	self.PerennialBranches = append(self.PerennialBranches, other.PerennialBranches...)
	if other.PerennialRegex.IsSome() {
		self.PerennialRegex = other.PerennialRegex
	}
	self.PrototypeBranches = append(self.PrototypeBranches, other.PrototypeBranches...)
	if other.PushHook.IsSome() {
		self.PushHook = other.PushHook
	}
	if other.ShipDeleteTrackingBranch.IsSome() {
		self.ShipDeleteTrackingBranch = other.ShipDeleteTrackingBranch
	}
	if other.SyncBeforeShip.IsSome() {
		self.SyncBeforeShip = other.SyncBeforeShip
	}
	if other.SyncFeatureStrategy.IsSome() {
		self.SyncFeatureStrategy = other.SyncFeatureStrategy
	}
	if other.SyncPerennialStrategy.IsSome() {
		self.SyncPerennialStrategy = other.SyncPerennialStrategy
	}
	if other.SyncUpstream.IsSome() {
		self.SyncUpstream = other.SyncUpstream
	}
}

func (self PartialConfig) ToUnvalidatedConfig(defaults UnvalidatedConfig) UnvalidatedConfig {
	return UnvalidatedConfig{
		Aliases:                  self.Aliases,
		ContributionBranches:     self.ContributionBranches,
		CreatePrototypeBranches:  self.CreatePrototypeBranches.GetOrElse(defaults.CreatePrototypeBranches),
		GitHubToken:              self.GitHubToken,
		GitLabToken:              self.GitLabToken,
		GitUserEmail:             self.GitUserEmail,
		GitUserName:              self.GitUserName,
		GiteaToken:               self.GiteaToken,
		HostingOriginHostname:    self.HostingOriginHostname,
		HostingPlatform:          self.HostingPlatform,
		Lineage:                  self.Lineage,
		MainBranch:               self.MainBranch,
		ObservedBranches:         self.ObservedBranches,
		Offline:                  self.Offline.GetOrElse(defaults.Offline),
		ParkedBranches:           self.ParkedBranches,
		PerennialBranches:        self.PerennialBranches,
		PerennialRegex:           self.PerennialRegex,
		PrototypeBranches:        self.PrototypeBranches,
		PushHook:                 self.PushHook.GetOrElse(defaults.PushHook),
		PushNewBranches:          self.PushNewBranches.GetOrElse(defaults.PushNewBranches),
		ShipDeleteTrackingBranch: self.ShipDeleteTrackingBranch.GetOrElse(defaults.ShipDeleteTrackingBranch),
		SyncBeforeShip:           self.SyncBeforeShip.GetOrElse(defaults.SyncBeforeShip),
		SyncFeatureStrategy:      self.SyncFeatureStrategy.GetOrElse(defaults.SyncFeatureStrategy),
		SyncPerennialStrategy:    self.SyncPerennialStrategy.GetOrElse(defaults.SyncPerennialStrategy),
		SyncUpstream:             self.SyncUpstream.GetOrElse(defaults.SyncUpstream),
	}
}
