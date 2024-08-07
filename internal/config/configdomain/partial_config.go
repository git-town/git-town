package configdomain

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks"
	"github.com/git-town/git-town/v15/internal/gohacks/mapstools"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
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
	return PartialConfig{
		Aliases:                  aliases,
		ContributionBranches:     gitdomain.ParseLocalBranchNames(snapshot[KeyContributionBranches]),
		CreatePrototypeBranches:  createPrototypeBranches,
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
		Offline:                  offline,
		ParkedBranches:           gitdomain.ParseLocalBranchNames(snapshot[KeyParkedBranches]),
		PerennialBranches:        gitdomain.ParseLocalBranchNames(snapshot[KeyPerennialBranches]),
		PerennialRegex:           ParsePerennialRegex(snapshot[KeyPerennialRegex]),
		PrototypeBranches:        gitdomain.ParseLocalBranchNames(snapshot[KeyPrototypeBranches]),
		PushHook:                 pushHook,
		PushNewBranches:          pushNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
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
		SyncFeatureStrategy:      other.SyncFeatureStrategy.Or(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.Or(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:    other.SyncPrototypeStrategy.Or(self.SyncPrototypeStrategy),
		SyncTags:                 other.SyncTags.Or(self.SyncTags),
		SyncUpstream:             other.SyncUpstream.Or(self.SyncUpstream),
	}
}

func (self PartialConfig) ToUnvalidatedConfig(defaults UnvalidatedConfig) UnvalidatedConfig {
	syncFeatureStrategy := self.SyncFeatureStrategy.GetOrElse(defaults.SyncFeatureStrategy)
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
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    self.SyncPerennialStrategy.GetOrElse(defaults.SyncPerennialStrategy),
		SyncPrototypeStrategy:    self.SyncPrototypeStrategy.GetOrElse(NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy)),
		SyncTags:                 self.SyncTags.GetOrElse(defaults.SyncTags),
		SyncUpstream:             self.SyncUpstream.GetOrElse(defaults.SyncUpstream),
	}
}
