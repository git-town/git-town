package configdomain

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/mapstools"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// PartialConfig contains configuration data as it is stored in one of the configuration sources for Git Town:
// - local Git metadata
// - global Git metadata
// - the configuration file
// - CLI arguments
//
// Any of these configuration data source can contain as much or as little configuration information as it wants.
// Hence, all fields here are optional.
type PartialConfig struct {
	Aliases                  Aliases
	AutoResolve              Option[AutoResolve]
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
	CodebergToken            Option[forgedomain.CodebergToken]
	ContributionRegex        Option[ContributionRegex]
	DevRemote                Option[gitdomain.Remote]
	DryRun                   Option[DryRun]
	FeatureRegex             Option[FeatureRegex]
	ForgeType                Option[forgedomain.ForgeType]
	GitHubConnectorType      Option[forgedomain.GitHubConnectorType]
	GitHubToken              Option[forgedomain.GitHubToken]
	GitLabConnectorType      Option[forgedomain.GitLabConnectorType]
	GitLabToken              Option[forgedomain.GitLabToken]
	GitUserEmail             Option[gitdomain.GitUserEmail]
	GitUserName              Option[gitdomain.GitUserName]
	GiteaToken               Option[forgedomain.GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	Lineage                  Lineage
	MainBranch               Option[gitdomain.LocalBranchName]
	NewBranchType            Option[NewBranchType]
	ObservedRegex            Option[ObservedRegex]
	Offline                  Option[Offline]
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	ProposalsShowLineage     Option[ProposalsShowLineage]
	PushHook                 Option[PushHook]
	ShareNewBranches         Option[ShareNewBranches]
	ShipDeleteTrackingBranch Option[ShipDeleteTrackingBranch]
	ShipStrategy             Option[ShipStrategy]
	SyncFeatureStrategy      Option[SyncFeatureStrategy]
	SyncPerennialStrategy    Option[SyncPerennialStrategy]
	SyncPrototypeStrategy    Option[SyncPrototypeStrategy]
	SyncTags                 Option[SyncTags]
	SyncUpstream             Option[SyncUpstream]
	UnknownBranchType        Option[UnknownBranchType]
	Verbose                  Option[Verbose]
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{
		Aliases: Aliases{},
		Lineage: NewLineage(),
	} //exhaustruct:ignore
}

// Merges the given PartialConfig into this configuration object.
func (self PartialConfig) Merge(other PartialConfig) PartialConfig {
	return PartialConfig{
		Aliases:                  mapstools.Merge(other.Aliases, self.Aliases),
		AutoResolve:              other.AutoResolve.Or(self.AutoResolve),
		BitbucketAppPassword:     other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:        other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchTypeOverrides:      other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		CodebergToken:            other.CodebergToken.Or(self.CodebergToken),
		ContributionRegex:        other.ContributionRegex.Or(self.ContributionRegex),
		DevRemote:                other.DevRemote.Or(self.DevRemote),
		DryRun:                   other.DryRun.Or(self.DryRun),
		FeatureRegex:             other.FeatureRegex.Or(self.FeatureRegex),
		ForgeType:                other.ForgeType.Or(self.ForgeType),
		GitHubConnectorType:      other.GitHubConnectorType.Or(self.GitHubConnectorType),
		GitHubToken:              other.GitHubToken.Or(self.GitHubToken),
		GitLabConnectorType:      other.GitLabConnectorType.Or(self.GitLabConnectorType),
		GitLabToken:              other.GitLabToken.Or(self.GitLabToken),
		GitUserEmail:             other.GitUserEmail.Or(self.GitUserEmail),
		GitUserName:              other.GitUserName.Or(self.GitUserName),
		GiteaToken:               other.GiteaToken.Or(self.GiteaToken),
		HostingOriginHostname:    other.HostingOriginHostname.Or(self.HostingOriginHostname),
		Lineage:                  other.Lineage.Merge(self.Lineage),
		MainBranch:               other.MainBranch.Or(self.MainBranch),
		NewBranchType:            other.NewBranchType.Or(self.NewBranchType),
		ObservedRegex:            other.ObservedRegex.Or(self.ObservedRegex),
		Offline:                  other.Offline.Or(self.Offline),
		PerennialBranches:        append(other.PerennialBranches, self.PerennialBranches...),
		PerennialRegex:           other.PerennialRegex.Or(self.PerennialRegex),
		ProposalsShowLineage:     other.ProposalsShowLineage.Or(self.ProposalsShowLineage),
		PushHook:                 other.PushHook.Or(self.PushHook),
		ShareNewBranches:         other.ShareNewBranches.Or(self.ShareNewBranches),
		ShipDeleteTrackingBranch: other.ShipDeleteTrackingBranch.Or(self.ShipDeleteTrackingBranch),
		ShipStrategy:             other.ShipStrategy.Or(self.ShipStrategy),
		SyncFeatureStrategy:      other.SyncFeatureStrategy.Or(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.Or(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:    other.SyncPrototypeStrategy.Or(self.SyncPrototypeStrategy),
		SyncTags:                 other.SyncTags.Or(self.SyncTags),
		SyncUpstream:             other.SyncUpstream.Or(self.SyncUpstream),
		UnknownBranchType:        other.UnknownBranchType.Or(self.UnknownBranchType),
		Verbose:                  other.Verbose.Or(self.Verbose),
	}
}

func (self PartialConfig) ToUnvalidatedConfig() UnvalidatedConfigData {
	return UnvalidatedConfigData{
		GitUserEmail: self.GitUserEmail,
		GitUserName:  self.GitUserName,
		MainBranch:   self.MainBranch,
	}
}
