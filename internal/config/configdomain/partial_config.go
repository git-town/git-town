package configdomain

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	AutoSync                 Option[AutoSync]
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
	ContributionRegex        Option[ContributionRegex]
	Detached                 Option[Detached]
	DevRemote                Option[gitdomain.Remote]
	DisplayTypes             Option[DisplayTypes]
	DryRun                   Option[DryRun]
	FeatureRegex             Option[FeatureRegex]
	ForgeType                Option[forgedomain.ForgeType]
	ForgejoToken             Option[forgedomain.ForgejoToken]
	GitHubConnectorType      Option[forgedomain.GitHubConnectorType]
	GitHubToken              Option[forgedomain.GitHubToken]
	GitHubUsername           Option[forgedomain.GitHubUsername]
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
	Order                    Option[Order]
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	ProposalsShowLineage     Option[forgedomain.ProposalsShowLineage]
	PushBranches             Option[PushBranches]
	PushHook                 Option[PushHook]
	ShareNewBranches         Option[ShareNewBranches]
	ShipDeleteTrackingBranch Option[ShipDeleteTrackingBranch]
	ShipStrategy             Option[ShipStrategy]
	Stash                    Option[Stash]
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
		AutoSync:                 other.AutoSync.Or(self.AutoSync),
		BitbucketAppPassword:     other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:        other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchTypeOverrides:      other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		ContributionRegex:        other.ContributionRegex.Or(self.ContributionRegex),
		Detached:                 other.Detached.Or(self.Detached),
		DevRemote:                other.DevRemote.Or(self.DevRemote),
		DisplayTypes:             other.DisplayTypes.Or(self.DisplayTypes),
		DryRun:                   other.DryRun.Or(self.DryRun),
		FeatureRegex:             other.FeatureRegex.Or(self.FeatureRegex),
		ForgeType:                other.ForgeType.Or(self.ForgeType),
		ForgejoToken:             other.ForgejoToken.Or(self.ForgejoToken),
		GitHubConnectorType:      other.GitHubConnectorType.Or(self.GitHubConnectorType),
		GitHubToken:              other.GitHubToken.Or(self.GitHubToken),
		GitHubUsername:           other.GitHubUsername.Or(self.GitHubUsername),
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
		Order:                    other.Order.Or(self.Order),
		PerennialBranches:        append(other.PerennialBranches, self.PerennialBranches...),
		PerennialRegex:           other.PerennialRegex.Or(self.PerennialRegex),
		ProposalsShowLineage:     other.ProposalsShowLineage.Or(self.ProposalsShowLineage),
		PushBranches:             other.PushBranches.Or(self.PushBranches),
		PushHook:                 other.PushHook.Or(self.PushHook),
		ShareNewBranches:         other.ShareNewBranches.Or(self.ShareNewBranches),
		ShipDeleteTrackingBranch: other.ShipDeleteTrackingBranch.Or(self.ShipDeleteTrackingBranch),
		ShipStrategy:             other.ShipStrategy.Or(self.ShipStrategy),
		Stash:                    other.Stash.Or(self.Stash),
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
