package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/git-town/git-town/v21/pkg/set"
)

// configuration settings that exist in both UnvalidatedConfig and ValidatedConfig
type NormalConfigData struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
	CodebergToken            Option[forgedomain.CodebergToken]
	ContributionRegex        Option[ContributionRegex]
	DevRemote                gitdomain.Remote
	DryRun                   DryRun // whether to only print the Git commands but not execute them
	FeatureRegex             Option[FeatureRegex]
	ForgeType                Option[forgedomain.ForgeType] // None = auto-detect
	GitHubConnectorType      Option[forgedomain.GitHubConnectorType]
	GitHubToken              Option[forgedomain.GitHubToken]
	GitLabConnectorType      Option[forgedomain.GitLabConnectorType]
	GitLabToken              Option[forgedomain.GitLabToken]
	GiteaToken               Option[forgedomain.GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	Lineage                  Lineage
	NewBranchType            Option[BranchType]
	ObservedRegex            Option[ObservedRegex]
	Offline                  Offline
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PushHook                 PushHook
	ShareNewBranches         ShareNewBranches
	ShipDeleteTrackingBranch ShipDeleteTrackingBranch
	ShipStrategy             ShipStrategy
	SyncFeatureStrategy      SyncFeatureStrategy
	SyncPerennialStrategy    SyncPerennialStrategy
	SyncPrototypeStrategy    SyncPrototypeStrategy
	SyncTags                 SyncTags
	SyncUpstream             SyncUpstream
	UnknownBranchType        BranchType
	Verbose                  Verbose
}

func (self *NormalConfigData) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *NormalConfigData) PartialBranchType(branch gitdomain.LocalBranchName) BranchType {
	// check the branch type overrides
	if branchTypeOverride, hasBranchTypeOverride := self.BranchTypeOverrides[branch]; hasBranchTypeOverride {
		return branchTypeOverride
	}
	// check the configured branch lists
	if slices.Contains(self.PerennialBranches, branch) {
		return BranchTypePerennialBranch
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
	// branch doesn't match any of the overrides --> unknown branch type
	return self.UnknownBranchType
}

func (self *NormalConfigData) PartialBranchesOfType(branchType BranchType) gitdomain.LocalBranchNames {
	matching := set.New[gitdomain.LocalBranchName]()
	switch branchType {
	case BranchTypeContributionBranch:
	case BranchTypeFeatureBranch:
	case BranchTypeMainBranch:
		// main branch is stored in ValidatedConfig
	case BranchTypeObservedBranch:
	case BranchTypeParkedBranch:
	case BranchTypePerennialBranch:
		matching.Add(self.PerennialBranches...)
	case BranchTypePrototypeBranch:
	}
	for key, value := range self.BranchTypeOverrides {
		if value == branchType {
			matching.Add(key)
		}
	}
	return matching.Values()
}

func DefaultNormalConfig() NormalConfigData {
	return NormalConfigData{
		Aliases:                  Aliases{},
		BitbucketAppPassword:     None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:        None[forgedomain.BitbucketUsername](),
		BranchTypeOverrides:      BranchTypeOverrides{},
		CodebergToken:            None[forgedomain.CodebergToken](),
		ContributionRegex:        None[ContributionRegex](),
		DevRemote:                gitdomain.RemoteOrigin,
		DryRun:                   false,
		FeatureRegex:             None[FeatureRegex](),
		ForgeType:                None[forgedomain.ForgeType](),
		GitHubConnectorType:      None[forgedomain.GitHubConnectorType](),
		GitHubToken:              None[forgedomain.GitHubToken](),
		GitLabConnectorType:      None[forgedomain.GitLabConnectorType](),
		GitLabToken:              None[forgedomain.GitLabToken](),
		GiteaToken:               None[forgedomain.GiteaToken](),
		HostingOriginHostname:    None[HostingOriginHostname](),
		Lineage:                  NewLineage(),
		NewBranchType:            None[BranchType](),
		ObservedRegex:            None[ObservedRegex](),
		Offline:                  false,
		PerennialBranches:        gitdomain.LocalBranchNames{},
		PerennialRegex:           None[PerennialRegex](),
		PushHook:                 true,
		ShareNewBranches:         ShareNewBranchesNone,
		ShipDeleteTrackingBranch: true,
		ShipStrategy:             ShipStrategyAPI,
		SyncFeatureStrategy:      SyncFeatureStrategyMerge,
		SyncPerennialStrategy:    SyncPerennialStrategyRebase,
		SyncPrototypeStrategy:    SyncPrototypeStrategyRebase,
		SyncTags:                 true,
		SyncUpstream:             true,
		UnknownBranchType:        BranchTypeFeatureBranch,
		Verbose:                  false,
	}
}
