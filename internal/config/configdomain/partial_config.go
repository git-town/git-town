package configdomain

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/mapstools"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
	CodebergToken            Option[forgedomain.CodebergToken]
	ContributionRegex        Option[ContributionRegex]
	DevRemote                Option[gitdomain.Remote]
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
	NewBranchType            Option[BranchType]
	ObservedRegex            Option[ObservedRegex]
	Offline                  Option[Offline]
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PushHook                 Option[PushHook]
	ShareNewBranches         Option[ShareNewBranches]
	ShipDeleteTrackingBranch Option[ShipDeleteTrackingBranch]
	ShipStrategy             Option[ShipStrategy]
	SyncFeatureStrategy      Option[SyncFeatureStrategy]
	SyncPerennialStrategy    Option[SyncPerennialStrategy]
	SyncPrototypeStrategy    Option[SyncPrototypeStrategy]
	SyncTags                 Option[SyncTags]
	SyncUpstream             Option[SyncUpstream]
	UnknownBranchType        Option[BranchType]
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
		BitbucketAppPassword:     other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:        other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchTypeOverrides:      other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		CodebergToken:            other.CodebergToken.Or(self.CodebergToken),
		ContributionRegex:        other.ContributionRegex.Or(self.ContributionRegex),
		DevRemote:                other.DevRemote.Or(self.DevRemote),
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
	}
}

func (self PartialConfig) ToNormalConfig(defaults NormalConfigData) NormalConfigData {
	syncFeatureStrategy := self.SyncFeatureStrategy.GetOrElse(defaults.SyncFeatureStrategy)
	return NormalConfigData{
		Aliases:                  self.Aliases,
		BitbucketAppPassword:     self.BitbucketAppPassword,
		BitbucketUsername:        self.BitbucketUsername,
		BranchTypeOverrides:      self.BranchTypeOverrides,
		CodebergToken:            self.CodebergToken,
		ContributionRegex:        self.ContributionRegex,
		DevRemote:                self.DevRemote.GetOrElse(defaults.DevRemote),
		FeatureRegex:             self.FeatureRegex,
		ForgeType:                self.ForgeType,
		GitHubConnectorType:      self.GitHubConnectorType,
		GitHubToken:              self.GitHubToken,
		GitLabConnectorType:      self.GitLabConnectorType,
		GitLabToken:              self.GitLabToken,
		GiteaToken:               self.GiteaToken,
		HostingOriginHostname:    self.HostingOriginHostname,
		Lineage:                  self.Lineage,
		NewBranchType:            self.NewBranchType.Or(defaults.NewBranchType),
		ObservedRegex:            self.ObservedRegex,
		Offline:                  self.Offline.GetOrElse(defaults.Offline),
		PerennialBranches:        self.PerennialBranches,
		PerennialRegex:           self.PerennialRegex,
		PushHook:                 self.PushHook.GetOrElse(defaults.PushHook),
		ShareNewBranches:         self.ShareNewBranches.GetOrElse(defaults.ShareNewBranches),
		ShipDeleteTrackingBranch: self.ShipDeleteTrackingBranch.GetOrElse(defaults.ShipDeleteTrackingBranch),
		ShipStrategy:             self.ShipStrategy.GetOrElse(defaults.ShipStrategy),
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    self.SyncPerennialStrategy.GetOrElse(defaults.SyncPerennialStrategy),
		SyncPrototypeStrategy:    self.SyncPrototypeStrategy.GetOrElse(NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy)),
		SyncTags:                 self.SyncTags.GetOrElse(defaults.SyncTags),
		SyncUpstream:             self.SyncUpstream.GetOrElse(defaults.SyncUpstream),
		UnknownBranchType:        self.UnknownBranchType.GetOrElse(BranchTypeFeatureBranch),
	}
}

func (self PartialConfig) ToUnvalidatedConfig() UnvalidatedConfigData {
	return UnvalidatedConfigData{
		GitUserEmail: self.GitUserEmail,
		GitUserName:  self.GitUserName,
		MainBranch:   self.MainBranch,
	}
}
