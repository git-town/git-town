package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/git-town/git-town/v17/pkg/set"
)

// configuration settings that exist in both UnvalidatedConfig and ValidatedConfig
type NormalConfigData struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[BitbucketAppPassword]
	BitbucketUsername        Option[BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
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

func (self *NormalConfigData) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *NormalConfigData) Online() Online {
	return self.Offline.ToOnline()
}

func (self *NormalConfigData) PartialBranchType(branch gitdomain.LocalBranchName) BranchType {
	// check the branch type overrides
	if branchTypeOverride, hasBranchTypeOverride := self.BranchTypeOverrides[branch]; hasBranchTypeOverride {
		return branchTypeOverride
	}
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

func (self *NormalConfigData) PartialBranchesOfType(branchType BranchType) gitdomain.LocalBranchNames {
	matching := set.New[gitdomain.LocalBranchName]()
	switch branchType {
	case BranchTypeContributionBranch:
		matching.Add(self.ContributionBranches...)
	case BranchTypeFeatureBranch:
	case BranchTypeMainBranch:
	case BranchTypeObservedBranch:
		matching.Add(self.ObservedBranches...)
	case BranchTypeParkedBranch:
		matching.Add(self.ParkedBranches...)
	case BranchTypePerennialBranch:
		matching.Add(self.PerennialBranches...)
	case BranchTypePrototypeBranch:
		matching.Add(self.PrototypeBranches...)
	}
	for key, value := range self.BranchTypeOverrides {
		if value == branchType {
			matching.Add(key)
		}
	}
	return matching.Values()
}

func (self *NormalConfigData) ShouldPushNewBranches() bool {
	return self.PushNewBranches.IsTrue()
}

func DefaultNormalConfig() NormalConfigData {
	return NormalConfigData{
		Aliases:                  Aliases{},
		BitbucketAppPassword:     None[BitbucketAppPassword](),
		BitbucketUsername:        None[BitbucketUsername](),
		BranchTypeOverrides:      BranchTypeOverrides{},
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
