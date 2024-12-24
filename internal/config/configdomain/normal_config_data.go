package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// configuration settings that exist in both UnvalidatedConfig and ValidatedConfig
type NormalConfigData struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[BitbucketAppPassword]
	BitbucketUsername        Option[BitbucketUsername]
	ContributionBranches     gitdomain.LocalBranchNames
	ContributionRegex        Option[ContributionRegex]
	DefaultBranchType        BranchType
	DevRemote                gitdomain.Remote
	FeatureRegex             Option[FeatureRegex]
	GitHubToken              Option[GitHubToken]
	GitLabToken              Option[GitLabToken]
	GiteaToken               Option[GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	HostingPlatform          Option[HostingPlatform] // Some = override by user, None = auto-detect
	Lineage                  Lineage
	NewBranchType            BranchType
	ObservedBranches         gitdomain.LocalBranchNames
	ObservedRegex            Option[ObservedRegex]
	Offline                  Offline
	ParkedBranches           gitdomain.LocalBranchNames
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PrototypeBranches        gitdomain.LocalBranchNames
	PushHook                 PushHook
	PushNewBranches          PushNewBranches
	ShipDeleteTrackingBranch ShipDeleteTrackingBranch
	ShipStrategy             ShipStrategy
	SyncFeatureStrategy      SyncFeatureStrategy
	SyncPerennialStrategy    SyncPerennialStrategy
	SyncPrototypeStrategy    SyncPrototypeStrategy
	SyncTags                 SyncTags
	SyncUpstream             SyncUpstream
}

func (self *NormalConfigData) IsOnline() bool {
	return self.Online().IsTrue()
}

// TODO: delete
func (self *NormalConfigData) IsParkedBranch(branch gitdomain.LocalBranchName) bool {
	return slices.Contains(self.ParkedBranches, branch)
}

// TODO: delete
func (self *NormalConfigData) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	if slices.Contains(self.PerennialBranches, branch) {
		return true
	}
	if perennialRegex, has := self.PerennialRegex.Get(); has {
		return perennialRegex.MatchesBranch(branch)
	}
	return false
}

// TODO: delete
func (self *NormalConfigData) IsPrototypeBranch(branch gitdomain.LocalBranchName) bool {
	if slices.Contains(self.PrototypeBranches, branch) {
		return true
	}
	return self.DefaultBranchType == BranchTypePrototypeBranch
}

func (self *NormalConfigData) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *NormalConfigData) Online() Online {
	return self.Offline.ToOnline()
}

func (self *NormalConfigData) PartialBranchType(branch gitdomain.LocalBranchName) BranchType {
	// check the configured branch lists
	if slices.Contains(self.ContributionBranches, branch) {
		return BranchTypeContributionBranch
	}
	if slices.Contains(self.ObservedBranches, branch) {
		return BranchTypeObservedBranch
	}
	if slices.Contains(self.ParkedBranches, branch) {
		return BranchTypeParkedBranch
	}
	if slices.Contains(self.PerennialBranches, branch) {
		return BranchTypePerennialBranch
	}
	if slices.Contains(self.PrototypeBranches, branch) {
		return BranchTypePrototypeBranch
	}
	// check if a regex matches
	if regex, has := self.ContributionRegex.Get(); has && regex.MatchesBranch(branch) {
		return BranchTypeContributionBranch
	}
	if regex, has := self.FeatureRegex.Get(); has && regex.MatchesBranch(branch) {
		return BranchTypeFeatureBranch
	}
	if regex, has := self.ObservedRegex.Get(); has && regex.MatchesBranch(branch) {
		return BranchTypeObservedBranch
	}
	if regex, has := self.PerennialRegex.Get(); has && regex.MatchesBranch(branch) {
		return BranchTypePerennialBranch
	}
	// branch doesn't match any of the overrides --> default branch type
	return self.DefaultBranchType
}

func (self *NormalConfigData) SetByKey(key Key, value string) {
	switch key {
	case KeyDevRemote:
		self.DevRemote = gitdomain.Remote(value)
	case KeyAliasAppend:
	case KeyAliasCompress:
	case KeyAliasContribute:
	case KeyAliasDelete:
	case KeyAliasDiffParent:
	case KeyAliasHack:
	case KeyAliasObserve:
	case KeyAliasPark:
	case KeyAliasPrepend:
	case KeyAliasPropose:
	case KeyAliasRename:
	case KeyAliasRepo:
	case KeyAliasSetParent:
	case KeyAliasShip:
	case KeyAliasSync:
	case KeyBitbucketAppPassword:
	case KeyBitbucketUsername:
	case KeyContributionBranches:
	case KeyContributionRegex:
	case KeyDefaultBranchType:
	case KeyDeprecatedAliasKill:
	case KeyDeprecatedAliasRenameBranch:
	case KeyDeprecatedCodeHostingDriver:
	case KeyDeprecatedCodeHostingOriginHostname:
	case KeyDeprecatedCodeHostingPlatform:
	case KeyDeprecatedCreatePrototypeBranches:
	case KeyDeprecatedMainBranchName:
	case KeyDeprecatedNewBranchPushFlag:
	case KeyDeprecatedPerennialBranchNames:
	case KeyDeprecatedPullBranchStrategy:
	case KeyDeprecatedPushVerify:
	case KeyDeprecatedShipDeleteRemoteBranch:
	case KeyDeprecatedSyncStrategy:
	case KeyFeatureRegex:
	case KeyGitUserEmail:
	case KeyGitUserName:
	case KeyGiteaToken:
	case KeyGithubToken:
	case KeyGitlabToken:
	case KeyHostingOriginHostname:
	case KeyHostingPlatform:
	case KeyMainBranch:
	case KeyNewBranchType:
	case KeyObservedBranches:
	case KeyObservedRegex:
	case KeyObsoleteSyncBeforeShip:
	case KeyOffline:
	case KeyParkedBranches:
	case KeyPerennialBranches:
	case KeyPerennialRegex:
	case KeyPrototypeBranches:
	case KeyPushHook:
	case KeyPushNewBranches:
	case KeyShipDeleteTrackingBranch:
	case KeyShipStrategy:
	case KeySyncFeatureStrategy:
	case KeySyncPerennialStrategy:
	case KeySyncPrototypeStrategy:
	case KeySyncTags:
	case KeySyncUpstream:
	}
}

func (self *NormalConfigData) ShouldPushNewBranches() bool {
	return self.PushNewBranches.IsTrue()
}

func DefaultNormalConfig() NormalConfigData {
	return NormalConfigData{
		Aliases:                  Aliases{},
		BitbucketAppPassword:     None[BitbucketAppPassword](),
		BitbucketUsername:        None[BitbucketUsername](),
		ContributionBranches:     gitdomain.LocalBranchNames{},
		ContributionRegex:        None[ContributionRegex](),
		DefaultBranchType:        BranchTypeFeatureBranch,
		DevRemote:                gitdomain.RemoteOrigin,
		FeatureRegex:             None[FeatureRegex](),
		GitHubToken:              None[GitHubToken](),
		GitLabToken:              None[GitLabToken](),
		GiteaToken:               None[GiteaToken](),
		HostingOriginHostname:    None[HostingOriginHostname](),
		HostingPlatform:          None[HostingPlatform](),
		Lineage:                  NewLineage(),
		NewBranchType:            BranchTypeFeatureBranch,
		ObservedBranches:         gitdomain.LocalBranchNames{},
		ObservedRegex:            None[ObservedRegex](),
		Offline:                  false,
		ParkedBranches:           gitdomain.LocalBranchNames{},
		PerennialBranches:        gitdomain.LocalBranchNames{},
		PerennialRegex:           None[PerennialRegex](),
		PrototypeBranches:        gitdomain.LocalBranchNames{},
		PushHook:                 true,
		PushNewBranches:          false,
		ShipDeleteTrackingBranch: true,
		ShipStrategy:             ShipStrategyAPI,
		SyncFeatureStrategy:      SyncFeatureStrategyMerge,
		SyncPerennialStrategy:    SyncPerennialStrategyRebase,
		SyncPrototypeStrategy:    SyncPrototypeStrategyRebase,
		SyncTags:                 true,
		SyncUpstream:             true,
	}
}
