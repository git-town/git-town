package configdomain

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/mapstools"
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
func (self PartialConfig) Merge(other PartialConfig) PartialConfig {
	return PartialConfig{
		Aliases:                  mapstools.Merge(other.Aliases, self.Aliases),
		ContributionBranches:     append(other.ContributionBranches, self.ContributionBranches...),
		CreatePrototypeBranches:  other.CreatePrototypeBranches.Or(self.CreatePrototypeBranches),
		GitHubToken:              other.GitHubToken.Or(self.GitHubToken),
		GitLabToken:              other.GitLabToken.Or(self.GitLabToken),
		GitUserEmail:             other.GitUserEmail.Or(self.GitUserEmail),
		GitUserName:              other.GitUserName.Or(self.GitUserName),
		GiteaToken:               other.GiteaToken.Or(self.GiteaToken),
		HostingOriginHostname:    other.HostingOriginHostname.Or(self.HostingOriginHostname),
		HostingPlatform:          other.HostingPlatform.Or(self.HostingPlatform),
		Lineage:                  other.Lineage.Merge(self.Lineage),
		MainBranch:               other.MainBranch.Or(self.MainBranch),
		ObservedBranches:         append(other.ObservedBranches, self.ObservedBranches...),
		Offline:                  other.Offline.Or(self.Offline),
		ParkedBranches:           append(other.ParkedBranches, self.ParkedBranches...),
		PerennialBranches:        append(other.PerennialBranches, self.PerennialBranches...),
		PerennialRegex:           other.PerennialRegex.Or(self.PerennialRegex),
		PrototypeBranches:        append(other.PrototypeBranches, self.PrototypeBranches...),
		PushHook:                 other.PushHook.Or(self.PushHook),
		PushNewBranches:          other.PushNewBranches.Or(self.PushNewBranches),
		ShipDeleteTrackingBranch: other.ShipDeleteTrackingBranch.Or(self.ShipDeleteTrackingBranch),
		SyncBeforeShip:           other.SyncBeforeShip.Or(self.SyncBeforeShip),
		SyncFeatureStrategy:      other.SyncFeatureStrategy.Or(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.Or(self.SyncPerennialStrategy),
		SyncUpstream:             other.SyncUpstream.Or(self.SyncUpstream),
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
