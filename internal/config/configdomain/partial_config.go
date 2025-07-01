package configdomain

import (
	"cmp"

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

func NewPartialConfigFromSnapshot(snapshot SingleSnapshot, updateOutdated bool, removeLocalConfigValue removeLocalConfigValueFunc) (PartialConfig, error) {
	aliases := snapshot.Aliases()
	branchTypeOverrides, err1 := NewBranchTypeOverridesInSnapshot(snapshot, removeLocalConfigValue)
	contributionRegex, err2 := ParseContributionRegex(snapshot[KeyContributionRegex])
	featureRegex, err3 := ParseFeatureRegex(snapshot[KeyFeatureRegex])
	forgeType, err4 := forgedomain.ParseForgeType(snapshot[KeyForgeType])
	githubConnectorType, err5 := forgedomain.ParseGitHubConnectorType(snapshot[KeyGitHubConnectorType])
	gitlabConnectorType, err6 := forgedomain.ParseGitLabConnectorType(snapshot[KeyGitLabConnectorType])
	lineage, err7 := NewLineageFromSnapshot(snapshot, updateOutdated, removeLocalConfigValue)
	newBranchType, err8 := ParseBranchType(snapshot[KeyNewBranchType])
	observedRegex, err9 := ParseObservedRegex(snapshot[KeyObservedRegex])
	offline, err10 := ParseOffline(snapshot[KeyOffline], KeyOffline)
	perennialRegex, err11 := ParsePerennialRegex(snapshot[KeyPerennialRegex])
	pushHook, err12 := ParsePushHook(snapshot[KeyPushHook], KeyPushHook)
	shareNewBranches, err13 := ParseShareNewBranches(snapshot[KeyShareNewBranches], KeyShareNewBranches)
	shipDeleteTrackingBranch, err14 := ParseShipDeleteTrackingBranch(snapshot[KeyShipDeleteTrackingBranch], KeyShipDeleteTrackingBranch)
	shipStrategy, err15 := ParseShipStrategy(snapshot[KeyShipStrategy])
	syncFeatureStrategy, err16 := ParseSyncFeatureStrategy(snapshot[KeySyncFeatureStrategy])
	syncPerennialStrategy, err17 := ParseSyncPerennialStrategy(snapshot[KeySyncPerennialStrategy])
	syncPrototypeStrategy, err18 := ParseSyncPrototypeStrategy(snapshot[KeySyncPrototypeStrategy])
	syncTags, err19 := ParseSyncTags(snapshot[KeySyncTags], KeySyncTags)
	syncUpstream, err20 := ParseSyncUpstream(snapshot[KeySyncUpstream], KeySyncUpstream)
	unknownBranchType, err21 := ParseBranchType(snapshot[KeyUnknownBranchType])
	return PartialConfig{
		Aliases:                  aliases,
		BitbucketAppPassword:     forgedomain.ParseBitbucketAppPassword(snapshot[KeyBitbucketAppPassword]),
		BitbucketUsername:        forgedomain.ParseBitbucketUsername(snapshot[KeyBitbucketUsername]),
		BranchTypeOverrides:      branchTypeOverrides,
		CodebergToken:            forgedomain.ParseCodebergToken(snapshot[KeyCodebergToken]),
		ContributionRegex:        contributionRegex,
		DevRemote:                gitdomain.NewRemote(snapshot[KeyDevRemote]),
		FeatureRegex:             featureRegex,
		ForgeType:                forgeType,
		GitHubConnectorType:      githubConnectorType,
		GitHubToken:              forgedomain.ParseGitHubToken(snapshot[KeyGitHubToken]),
		GitLabConnectorType:      gitlabConnectorType,
		GitLabToken:              forgedomain.ParseGitLabToken(snapshot[KeyGitLabToken]),
		GitUserEmail:             gitdomain.ParseGitUserEmail(snapshot[KeyGitUserEmail]),
		GitUserName:              gitdomain.ParseGitUserName(snapshot[KeyGitUserName]),
		GiteaToken:               forgedomain.ParseGiteaToken(snapshot[KeyGiteaToken]),
		HostingOriginHostname:    ParseHostingOriginHostname(snapshot[KeyHostingOriginHostname]),
		Lineage:                  lineage,
		MainBranch:               gitdomain.NewLocalBranchNameOption(snapshot[KeyMainBranch]),
		NewBranchType:            newBranchType,
		ObservedRegex:            observedRegex,
		Offline:                  offline,
		PerennialBranches:        gitdomain.ParseLocalBranchNames(snapshot[KeyPerennialBranches]),
		PerennialRegex:           perennialRegex,
		PushHook:                 pushHook,
		ShareNewBranches:         shareNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
		UnknownBranchType:        unknownBranchType,
	}, cmp.Or(err1, err2, err3, err4, err5, err6, err7, err8, err9, err10, err11, err12, err13, err14, err15, err16, err17, err18, err19, err20, err21)
}

// a function that deletes the local Git configuration value with the given key
type removeLocalConfigValueFunc func(Key) error

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
