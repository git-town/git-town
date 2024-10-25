package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// configuration settings that exist in both UnvalidatedConfig and ValidatedConfig
type NormalConfigData struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[BitbucketAppPassword]
	BitbucketUsername        Option[BitbucketUsername]
	ContributionBranches     gitdomain.LocalBranchNames
	ContributionRegex        Option[ContributionRegex]
	CreatePrototypeBranches  CreatePrototypeBranches
	DefaultBranchType        DefaultBranchType
	FeatureRegex             Option[FeatureRegex]
	GitHubToken              Option[GitHubToken]
	GitLabToken              Option[GitLabToken]
	GiteaToken               Option[GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	HostingPlatform          Option[HostingPlatform] // Some = override by user, None = auto-detect
	Lineage                  Lineage
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

// ContainsLineage indicates whether this configuration contains any lineage entries.
func (self *NormalConfigData) ContainsLineage() bool {
	return self.Lineage.Len() > 0
}

func (self *NormalConfigData) IsOnline() bool {
	return self.Online().IsTrue()
}

func (self *NormalConfigData) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	if slices.Contains(self.PerennialBranches, branch) {
		return true
	}
	if perennialRegex, has := self.PerennialRegex.Get(); has {
		return perennialRegex.MatchesBranch(branch)
	}
	return false
}

func (self *NormalConfigData) MatchesContributionRegex(branch gitdomain.LocalBranchName) bool {
	if contributionRegex, has := self.ContributionRegex.Get(); has {
		return contributionRegex.MatchesBranch(branch)
	}
	return false
}

func (self *NormalConfigData) MatchesFeatureBranchRegex(branch gitdomain.LocalBranchName) bool {
	if featureRegex, has := self.FeatureRegex.Get(); has {
		return featureRegex.MatchesBranch(branch)
	}
	return false
}

func (self *NormalConfigData) MatchesObservedRegex(branch gitdomain.LocalBranchName) bool {
	if observedRegex, has := self.ObservedRegex.Get(); has {
		return observedRegex.MatchesBranch(branch)
	}
	return false
}

func (self *NormalConfigData) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *NormalConfigData) Online() Online {
	return self.Offline.ToOnline()
}

func (self *NormalConfigData) PartialBranchType(branch gitdomain.LocalBranchName) BranchType {
	if self.IsPerennialBranch(branch) {
		return BranchTypePerennialBranch
	}
	if slices.Contains(self.ContributionBranches, branch) {
		return BranchTypeContributionBranch
	}
	if slices.Contains(self.ObservedBranches, branch) {
		return BranchTypeObservedBranch
	}
	if slices.Contains(self.ParkedBranches, branch) {
		return BranchTypeParkedBranch
	}
	if slices.Contains(self.PrototypeBranches, branch) {
		return BranchTypePrototypeBranch
	}
	if self.MatchesFeatureBranchRegex(branch) {
		return BranchTypeFeatureBranch
	}
	if self.MatchesContributionRegex(branch) {
		return BranchTypeContributionBranch
	}
	if self.MatchesObservedRegex(branch) {
		return BranchTypeObservedBranch
	}
	return self.DefaultBranchType.BranchType
}

func (self *NormalConfigData) ShouldPushNewBranches() bool {
	return self.PushNewBranches.IsTrue()
}

func DefaultNormalConfig() NormalConfigData {
	return NormalConfigData{
		Aliases:                  Aliases{},
		BitbucketAppPassword:     None[BitbucketAppPassword](),
		BitbucketUsername:        None[BitbucketUsername](),
		ContributionBranches:     gitdomain.NewLocalBranchNames(),
		ContributionRegex:        None[ContributionRegex](),
		CreatePrototypeBranches:  false,
		DefaultBranchType:        DefaultBranchType{BranchType: BranchTypeFeatureBranch},
		FeatureRegex:             None[FeatureRegex](),
		GitHubToken:              None[GitHubToken](),
		GitLabToken:              None[GitLabToken](),
		GiteaToken:               None[GiteaToken](),
		HostingOriginHostname:    None[HostingOriginHostname](),
		HostingPlatform:          None[HostingPlatform](),
		Lineage:                  NewLineage(),
		ObservedBranches:         gitdomain.NewLocalBranchNames(),
		ObservedRegex:            None[ObservedRegex](),
		Offline:                  false,
		ParkedBranches:           gitdomain.NewLocalBranchNames(),
		PerennialBranches:        gitdomain.NewLocalBranchNames(),
		PerennialRegex:           None[PerennialRegex](),
		PrototypeBranches:        gitdomain.NewLocalBranchNames(),
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
