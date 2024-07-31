package configdomain

import (
	"strings"

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
	SyncPrototypeStrategy    Option[SyncPrototypeStrategy]
	SyncUpstream             Option[SyncUpstream]
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{
		Aliases: Aliases{},
	} //exhaustruct:ignore
}

// a function that deletes the local Git configuration value with the given key
type removeLocalConfigValueFunc func(Key) error

// Note: this exists here and not as a method of PartialConfig to avoid circular dependencies
func (self *PartialConfig) AddValue(key Key, value string, removeLocalConfigValue removeLocalConfigValueFunc) error {
	if strings.HasPrefix(key.String(), LineageKeyPrefix) {
		if childStr, isLineage := key.IsLineage().Get(); isLineage {
			if childStr == "" {
				// empty lineage entries are invalid --> delete it
				return removeLocalConfigValue(key)
			}
			child := gitdomain.NewLocalBranchName(childStr)
			value = strings.TrimSpace(value)
			if value == "" {
				// empty lineage entries are invalid --> delete it
				return removeLocalConfigValue(key)
			}
			parent := gitdomain.NewLocalBranchName(value)
			self.Lineage.Add(child, parent)
			return nil
		}
	}
	var err error
	switch key {
	case KeyAliasAppend:
		self.Aliases[AliasableCommandAppend] = value
	case KeyAliasCompress:
		self.Aliases[AliasableCommandCompress] = value
	case KeyAliasContribute:
		self.Aliases[AliasableCommandContribute] = value
	case KeyAliasDiffParent:
		self.Aliases[AliasableCommandDiffParent] = value
	case KeyAliasHack:
		self.Aliases[AliasableCommandHack] = value
	case KeyAliasKill:
		self.Aliases[AliasableCommandKill] = value
	case KeyAliasObserve:
		self.Aliases[AliasableCommandObserve] = value
	case KeyAliasPark:
		self.Aliases[AliasableCommandPark] = value
	case KeyAliasPrepend:
		self.Aliases[AliasableCommandPrepend] = value
	case KeyAliasPropose:
		self.Aliases[AliasableCommandPropose] = value
	case KeyAliasRenameBranch:
		self.Aliases[AliasableCommandRenameBranch] = value
	case KeyAliasRepo:
		self.Aliases[AliasableCommandRepo] = value
	case KeyAliasSetParent:
		self.Aliases[AliasableCommandSetParent] = value
	case KeyAliasShip:
		self.Aliases[AliasableCommandShip] = value
	case KeyAliasSync:
		self.Aliases[AliasableCommandSync] = value
	case KeyContributionBranches:
		self.ContributionBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyCreatePrototypeBranches:
		self.CreatePrototypeBranches, err = ParseCreatePrototypeBranchesOpt(value, KeyPrototypeBranches.String())
	case KeyHostingOriginHostname:
		self.HostingOriginHostname = NewHostingOriginHostnameOption(value)
	case KeyHostingPlatform:
		self.HostingPlatform, err = NewHostingPlatformOption(value)
	case KeyGiteaToken:
		self.GiteaToken = NewGiteaTokenOption(value)
	case KeyGithubToken:
		self.GitHubToken = NewGitHubTokenOption(value)
	case KeyGitlabToken:
		self.GitLabToken = NewGitLabTokenOption(value)
	case KeyGitUserEmail:
		self.GitUserEmail = NewGitUserEmailOption(value)
	case KeyGitUserName:
		self.GitUserName = NewGitUserNameOption(value)
	case KeyMainBranch:
		self.MainBranch = gitdomain.NewLocalBranchNameOption(value)
	case KeyObservedBranches:
		self.ObservedBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyOffline:
		self.Offline, err = NewOfflineOption(value, KeyOffline.String())
	case KeyParkedBranches:
		self.ParkedBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyPerennialBranches:
		self.PerennialBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyPerennialRegex:
		self.PerennialRegex = NewPerennialRegexOption(value)
	case KeyPrototypeBranches:
		self.PrototypeBranches = gitdomain.ParseLocalBranchNames(value)
	case KeyPushHook:
		var pushHook PushHook
		pushHook, err = NewPushHook(value, KeyPushHook.String())
		self.PushHook = Some(pushHook)
	case KeyPushNewBranches:
		self.PushNewBranches, err = ParsePushNewBranchesOption(value, KeyPushNewBranches.String())
	case KeyShipDeleteTrackingBranch:
		self.ShipDeleteTrackingBranch, err = ParseShipDeleteTrackingBranchOption(value, KeyShipDeleteTrackingBranch.String())
	case KeySyncFeatureStrategy:
		self.SyncFeatureStrategy, err = NewSyncFeatureStrategyOption(value)
	case KeySyncPerennialStrategy:
		self.SyncPerennialStrategy, err = NewSyncPerennialStrategyOption(value)
	case KeySyncPrototypeStrategy:
		self.SyncPrototypeStrategy, err = NewSyncPrototypeStrategyOption(value)
	case KeySyncUpstream:
		self.SyncUpstream, err = ParseSyncUpstreamOption(value, KeySyncUpstream.String())
	case KeyDeprecatedCodeHostingDriver,
		KeyDeprecatedCodeHostingOriginHostname,
		KeyDeprecatedCodeHostingPlatform,
		KeyDeprecatedMainBranchName,
		KeyDeprecatedNewBranchPushFlag,
		KeyDeprecatedPerennialBranchNames,
		KeyDeprecatedPullBranchStrategy,
		KeyDeprecatedPushVerify,
		KeyDeprecatedShipDeleteRemoteBranch,
		KeyDeprecatedSyncStrategy:
		// deprecated keys were handled before this is reached, they are listed here to check that the switch statement contains all keys
	}
	return err
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
		SyncFeatureStrategy:      other.SyncFeatureStrategy.Or(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.Or(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:    other.SyncPrototypeStrategy.Or(self.SyncPrototypeStrategy),
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
		SyncUpstream:             self.SyncUpstream.GetOrElse(defaults.SyncUpstream),
	}
}
