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
	Aliases                     Aliases
	AutoResolve                 Option[AutoResolve]
	AutoSync                    Option[AutoSync]
	BitbucketAppPassword        Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername           Option[forgedomain.BitbucketUsername]
	BranchPrefix                Option[BranchPrefix]
	BranchTypeOverrides         BranchTypeOverrides
	Browser                     Option[Browser]
	ContributionRegex           Option[ContributionRegex]
	Detached                    Option[Detached]
	DevRemote                   Option[gitdomain.Remote]
	DisplayTypes                Option[DisplayTypes]
	DryRun                      Option[DryRun]
	FeatureRegex                Option[FeatureRegex]
	ForgeType                   Option[forgedomain.ForgeType]
	ForgejoToken                Option[forgedomain.ForgejoToken]
	GitUserEmail                Option[gitdomain.GitUserEmail]
	GitUserName                 Option[gitdomain.GitUserName]
	GiteaToken                  Option[forgedomain.GiteaToken]
	GithubConnectorType         Option[forgedomain.GithubConnectorType]
	GithubToken                 Option[forgedomain.GithubToken]
	GitlabConnectorType         Option[forgedomain.GitlabConnectorType]
	GitlabToken                 Option[forgedomain.GitlabToken]
	HostingOriginHostname       Option[HostingOriginHostname]
	IgnoreUncommitted           Option[IgnoreUncommitted]
	Lineage                     Lineage
	MainBranch                  Option[gitdomain.LocalBranchName]
	NewBranchType               Option[NewBranchType]
	ObservedRegex               Option[ObservedRegex]
	Offline                     Option[Offline]
	Order                       Option[Order]
	PerennialBranches           gitdomain.LocalBranchNames
	PerennialRegex              Option[PerennialRegex]
	ProposalBreadcrumb          Option[ProposalBreadcrumb]
	ProposalBreadcrumbDirection Option[ProposalBreadcrumbDirection]
	PushBranches                Option[PushBranches]
	PushHook                    Option[PushHook]
	ShareNewBranches            Option[ShareNewBranches]
	ShipDeleteTrackingBranch    Option[ShipDeleteTrackingBranch]
	ShipStrategy                Option[ShipStrategy]
	Stash                       Option[Stash]
	SyncFeatureStrategy         Option[SyncFeatureStrategy]
	SyncPerennialStrategy       Option[SyncPerennialStrategy]
	SyncPrototypeStrategy       Option[SyncPrototypeStrategy]
	SyncTags                    Option[SyncTags]
	SyncUpstream                Option[SyncUpstream]
	UnknownBranchType           Option[UnknownBranchType]
	Verbose                     Option[Verbose]
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{
		Aliases: Aliases{},
		Lineage: NewLineage(),
	} //exhaustruct:ignore
}

// Merge combines the data of the given PartialConfig with this PartialConfig,
// favoring the data of the given PartialConfig.
func (self PartialConfig) Merge(other PartialConfig) PartialConfig {
	return PartialConfig{
		Aliases:                     mapstools.Merge(other.Aliases, self.Aliases),
		AutoResolve:                 other.AutoResolve.Or(self.AutoResolve),
		AutoSync:                    other.AutoSync.Or(self.AutoSync),
		BitbucketAppPassword:        other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:           other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchPrefix:                other.BranchPrefix.Or(self.BranchPrefix),
		BranchTypeOverrides:         other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		Browser:                     other.Browser.Or(self.Browser),
		ContributionRegex:           other.ContributionRegex.Or(self.ContributionRegex),
		Detached:                    other.Detached.Or(self.Detached),
		DevRemote:                   other.DevRemote.Or(self.DevRemote),
		DisplayTypes:                other.DisplayTypes.Or(self.DisplayTypes),
		DryRun:                      other.DryRun.Or(self.DryRun),
		FeatureRegex:                other.FeatureRegex.Or(self.FeatureRegex),
		ForgeType:                   other.ForgeType.Or(self.ForgeType),
		ForgejoToken:                other.ForgejoToken.Or(self.ForgejoToken),
		GitUserEmail:                other.GitUserEmail.Or(self.GitUserEmail),
		GitUserName:                 other.GitUserName.Or(self.GitUserName),
		GiteaToken:                  other.GiteaToken.Or(self.GiteaToken),
		GithubConnectorType:         other.GithubConnectorType.Or(self.GithubConnectorType),
		GithubToken:                 other.GithubToken.Or(self.GithubToken),
		GitlabConnectorType:         other.GitlabConnectorType.Or(self.GitlabConnectorType),
		GitlabToken:                 other.GitlabToken.Or(self.GitlabToken),
		HostingOriginHostname:       other.HostingOriginHostname.Or(self.HostingOriginHostname),
		IgnoreUncommitted:           other.IgnoreUncommitted.Or(self.IgnoreUncommitted),
		Lineage:                     other.Lineage.Merge(self.Lineage),
		MainBranch:                  other.MainBranch.Or(self.MainBranch),
		NewBranchType:               other.NewBranchType.Or(self.NewBranchType),
		ObservedRegex:               other.ObservedRegex.Or(self.ObservedRegex),
		Offline:                     other.Offline.Or(self.Offline),
		Order:                       other.Order.Or(self.Order),
		PerennialBranches:           append(other.PerennialBranches, self.PerennialBranches...),
		PerennialRegex:              other.PerennialRegex.Or(self.PerennialRegex),
		ProposalBreadcrumb:          other.ProposalBreadcrumb.Or(self.ProposalBreadcrumb),
		ProposalBreadcrumbDirection: other.ProposalBreadcrumbDirection.Or(self.ProposalBreadcrumbDirection),
		PushBranches:                other.PushBranches.Or(self.PushBranches),
		PushHook:                    other.PushHook.Or(self.PushHook),
		ShareNewBranches:            other.ShareNewBranches.Or(self.ShareNewBranches),
		ShipDeleteTrackingBranch:    other.ShipDeleteTrackingBranch.Or(self.ShipDeleteTrackingBranch),
		ShipStrategy:                other.ShipStrategy.Or(self.ShipStrategy),
		Stash:                       other.Stash.Or(self.Stash),
		SyncFeatureStrategy:         other.SyncFeatureStrategy.Or(self.SyncFeatureStrategy),
		SyncPerennialStrategy:       other.SyncPerennialStrategy.Or(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:       other.SyncPrototypeStrategy.Or(self.SyncPrototypeStrategy),
		SyncTags:                    other.SyncTags.Or(self.SyncTags),
		SyncUpstream:                other.SyncUpstream.Or(self.SyncUpstream),
		UnknownBranchType:           other.UnknownBranchType.Or(self.UnknownBranchType),
		Verbose:                     other.Verbose.Or(self.Verbose),
	}
}

func (self PartialConfig) ToUnvalidatedConfig() UnvalidatedConfigData {
	return UnvalidatedConfigData{
		GitUserEmail: self.GitUserEmail,
		GitUserName:  self.GitUserName,
		MainBranch:   self.MainBranch,
	}
}
