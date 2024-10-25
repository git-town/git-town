package configdomain

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks"
	"github.com/git-town/git-town/v16/internal/gohacks/mapstools"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[BitbucketAppPassword]
	BitbucketUsername        Option[BitbucketUsername]
	ContributionBranches     gitdomain.LocalBranchNames
	ContributionRegex        Option[ContributionRegex]
	CreatePrototypeBranches  Option[CreatePrototypeBranches]
	DefaultBranchType        Option[DefaultBranchType]
	FeatureRegex             Option[FeatureRegex]
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
	ObservedRegex            Option[ObservedRegex]
	Offline                  Option[Offline]
	ParkedBranches           gitdomain.LocalBranchNames
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PrototypeBranches        gitdomain.LocalBranchNames
	PushHook                 Option[PushHook]
	PushNewBranches          Option[PushNewBranches]
	ShipDeleteTrackingBranch Option[ShipDeleteTrackingBranch]
	ShipStrategy             Option[ShipStrategy]
	SyncFeatureStrategy      Option[SyncFeatureStrategy]
	SyncPerennialStrategy    Option[SyncPerennialStrategy]
	SyncPrototypeStrategy    Option[SyncPrototypeStrategy]
	SyncTags                 Option[SyncTags]
	SyncUpstream             Option[SyncUpstream]
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{
		Aliases: Aliases{},
	} //exhaustruct:ignore
}

func NewPartialConfigFromSnapshot(snapshot SingleSnapshot, updateOutdated bool, removeLocalConfigValue removeLocalConfigValueFunc) (PartialConfig, error) {
	ec := gohacks.ErrorCollector{}
	aliases := snapshot.Aliases()
	createPrototypeBranches, err := ParseCreatePrototypeBranches(snapshot[KeyCreatePrototypeBranches], KeyCreatePrototypeBranches.String())
	ec.Check(err)
	hostingPlatform, err := ParseHostingPlatform(snapshot[KeyHostingPlatform])
	ec.Check(err)
	offline, err := ParseOffline(snapshot[KeyOffline], KeyOffline.String())
	ec.Check(err)
	pushHook, err := ParsePushHook(snapshot[KeyPushHook], KeyPushHook.String())
	ec.Check(err)
	pushNewBranches, err := ParsePushNewBranches(snapshot[KeyPushNewBranches], KeyPushNewBranches.String())
	ec.Check(err)
	shipDeleteTrackingBranch, err := ParseShipDeleteTrackingBranch(snapshot[KeyShipDeleteTrackingBranch], KeyShipDeleteTrackingBranch.String())
	ec.Check(err)
	shipStrategy, err := ParseShipStrategy(snapshot[KeyShipStrategy])
	ec.Check(err)
	syncFeatureStrategy, err := ParseSyncFeatureStrategy(snapshot[KeySyncFeatureStrategy])
	ec.Check(err)
	syncPerennialStrategy, err := ParseSyncPerennialStrategy(snapshot[KeySyncPerennialStrategy])
	ec.Check(err)
	syncPrototypeStrategy, err := ParseSyncPrototypeStrategy(snapshot[KeySyncPrototypeStrategy])
	ec.Check(err)
	syncTags, err := ParseSyncTags(snapshot[KeySyncTags], KeySyncTags.String())
	ec.Check(err)
	syncUpstream, err := ParseSyncUpstream(snapshot[KeySyncUpstream], KeySyncUpstream.String())
	ec.Check(err)
	lineage, err := NewLineageFromSnapshot(snapshot, updateOutdated, removeLocalConfigValue)
	ec.Check(err)
	perennialRegex, err := ParsePerennialRegex(snapshot[KeyPerennialRegex])
	ec.Check(err)
	defaultBranchType, err := ParseDefaultBranchType(snapshot[KeyDefaultBranchType])
	ec.Check(err)
	featureRegex, err := ParseFeatureRegex(snapshot[KeyFeatureRegex])
	ec.Check(err)
	contributionRegex, err := ParseContributionRegex(snapshot[KeyContributionRegex])
	ec.Check(err)
	observedRegex, err := ParseObservedRegex(snapshot[KeyObservedRegex])
	ec.Check(err)
	return PartialConfig{
		Aliases:                  aliases,
		BitbucketAppPassword:     ParseBitbucketAppPassword(snapshot[KeyBitbucketAppPassword]),
		BitbucketUsername:        ParseBitbucketUsername(snapshot[KeyBitbucketUsername]),
		ContributionBranches:     gitdomain.ParseLocalBranchNames(snapshot[KeyContributionBranches]),
		ContributionRegex:        contributionRegex,
		CreatePrototypeBranches:  createPrototypeBranches,
		DefaultBranchType:        defaultBranchType,
		FeatureRegex:             featureRegex,
		GitHubToken:              ParseGitHubToken(snapshot[KeyGithubToken]),
		GitLabToken:              ParseGitLabToken(snapshot[KeyGitlabToken]),
		GitUserEmail:             ParseGitUserEmail(snapshot[KeyGitUserEmail]),
		GitUserName:              ParseGitUserName(snapshot[KeyGitUserName]),
		GiteaToken:               ParseGiteaToken(snapshot[KeyGiteaToken]),
		HostingOriginHostname:    ParseHostingOriginHostname(snapshot[KeyHostingOriginHostname]),
		HostingPlatform:          hostingPlatform,
		Lineage:                  lineage,
		MainBranch:               gitdomain.NewLocalBranchNameOption(snapshot[KeyMainBranch]),
		ObservedBranches:         gitdomain.ParseLocalBranchNames(snapshot[KeyObservedBranches]),
		ObservedRegex:            observedRegex,
		Offline:                  offline,
		ParkedBranches:           gitdomain.ParseLocalBranchNames(snapshot[KeyParkedBranches]),
		PerennialBranches:        gitdomain.ParseLocalBranchNames(snapshot[KeyPerennialBranches]),
		PerennialRegex:           perennialRegex,
		PrototypeBranches:        gitdomain.ParseLocalBranchNames(snapshot[KeyPrototypeBranches]),
		PushHook:                 pushHook,
		PushNewBranches:          pushNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
	}, ec.Err
}

// a function that deletes the local Git configuration value with the given key
type removeLocalConfigValueFunc func(Key) error

// Merges the given PartialConfig into this configuration object.
func (self PartialConfig) Merge(other PartialConfig) PartialConfig {
	return PartialConfig{
		Aliases:                  mapstools.Merge(other.Aliases, self.Aliases),
		BitbucketAppPassword:     other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:        other.BitbucketUsername.Or(self.BitbucketUsername),
		ContributionBranches:     append(other.ContributionBranches, self.ContributionBranches...),
		ContributionRegex:        other.ContributionRegex.Or(self.ContributionRegex),
		CreatePrototypeBranches:  other.CreatePrototypeBranches.Or(self.CreatePrototypeBranches),
		DefaultBranchType:        other.DefaultBranchType.Or(self.DefaultBranchType),
		FeatureRegex:             other.FeatureRegex.Or(self.FeatureRegex),
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
		ObservedRegex:            other.ObservedRegex.Or(self.ObservedRegex),
		Offline:                  other.Offline.Or(self.Offline),
		ParkedBranches:           append(other.ParkedBranches, self.ParkedBranches...),
		PerennialBranches:        append(other.PerennialBranches, self.PerennialBranches...),
		PerennialRegex:           other.PerennialRegex.Or(self.PerennialRegex),
		PrototypeBranches:        append(other.PrototypeBranches, self.PrototypeBranches...),
		PushHook:                 other.PushHook.Or(self.PushHook),
		PushNewBranches:          other.PushNewBranches.Or(self.PushNewBranches),
		ShipDeleteTrackingBranch: other.ShipDeleteTrackingBranch.Or(self.ShipDeleteTrackingBranch),
		ShipStrategy:             other.ShipStrategy.Or(self.ShipStrategy),
		SyncFeatureStrategy:      other.SyncFeatureStrategy.Or(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.Or(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:    other.SyncPrototypeStrategy.Or(self.SyncPrototypeStrategy),
		SyncTags:                 other.SyncTags.Or(self.SyncTags),
		SyncUpstream:             other.SyncUpstream.Or(self.SyncUpstream),
	}
}

func (self PartialConfig) ToNormalConfig(defaults NormalConfigData) NormalConfigData {
	syncFeatureStrategy := self.SyncFeatureStrategy.GetOrElse(defaults.SyncFeatureStrategy)
	return NormalConfigData{
		Aliases:                  self.Aliases,
		BitbucketAppPassword:     self.BitbucketAppPassword,
		BitbucketUsername:        self.BitbucketUsername,
		ContributionBranches:     self.ContributionBranches,
		ContributionRegex:        self.ContributionRegex,
		CreatePrototypeBranches:  self.CreatePrototypeBranches.GetOrElse(defaults.CreatePrototypeBranches),
		DefaultBranchType:        self.DefaultBranchType.GetOrElse(DefaultBranchType{BranchTypeFeatureBranch}),
		FeatureRegex:             self.FeatureRegex,
		GitHubToken:              self.GitHubToken,
		GitLabToken:              self.GitLabToken,
		GiteaToken:               self.GiteaToken,
		HostingOriginHostname:    self.HostingOriginHostname,
		HostingPlatform:          self.HostingPlatform,
		Lineage:                  self.Lineage,
		ObservedBranches:         self.ObservedBranches,
		ObservedRegex:            self.ObservedRegex,
		Offline:                  self.Offline.GetOrElse(defaults.Offline),
		ParkedBranches:           self.ParkedBranches,
		PerennialBranches:        self.PerennialBranches,
		PerennialRegex:           self.PerennialRegex,
		PrototypeBranches:        self.PrototypeBranches,
		PushHook:                 self.PushHook.GetOrElse(defaults.PushHook),
		PushNewBranches:          self.PushNewBranches.GetOrElse(defaults.PushNewBranches),
		ShipDeleteTrackingBranch: self.ShipDeleteTrackingBranch.GetOrElse(defaults.ShipDeleteTrackingBranch),
		ShipStrategy:             self.ShipStrategy.GetOrElse(defaults.ShipStrategy),
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    self.SyncPerennialStrategy.GetOrElse(defaults.SyncPerennialStrategy),
		SyncPrototypeStrategy:    self.SyncPrototypeStrategy.GetOrElse(NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy)),
		SyncTags:                 self.SyncTags.GetOrElse(defaults.SyncTags),
		SyncUpstream:             self.SyncUpstream.GetOrElse(defaults.SyncUpstream),
	}
}

func (self PartialConfig) ToUnvalidatedConfig() UnvalidatedConfigData {
	return UnvalidatedConfigData{
		GitUserEmail: self.GitUserEmail,
		GitUserName:  self.GitUserName,
		MainBranch:   self.MainBranch,
	}
}
