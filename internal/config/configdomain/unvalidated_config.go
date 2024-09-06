package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// UnvalidatedConfig is the Git Town configuration as read from disk.
// It might be lacking essential information in case Git metadata and config files don't contain it.
// If you need this information, validate it into a ValidatedConfig.
type UnvalidatedConfig struct {
	Aliases                  Aliases
	ContributionBranches     gitdomain.LocalBranchNames
	CreatePrototypeBranches  CreatePrototypeBranches
	DefaultBranchType        DefaultBranchType
	FeatureRegex             Option[FeatureRegex]
	GitHubToken              Option[GitHubToken]
	GitLabToken              Option[GitLabToken]
	GitUserEmail             Option[GitUserEmail]
	GitUserName              Option[GitUserName]
	GiteaToken               Option[GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	HostingPlatform          Option[HostingPlatform] // Some = override by user, None = auto-detect
	Lineage                  Lineage
	MainBranch               Option[gitdomain.LocalBranchName]
	ObservedBranches         gitdomain.LocalBranchNames
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

// indicates the branch type of the given branch
func (self *UnvalidatedConfig) BranchType(branch gitdomain.LocalBranchName) BranchType {
	if self.IsMainBranch(branch) {
		return BranchTypeMainBranch
	}
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
	return self.DefaultBranchType.BranchType
}

func (self *UnvalidatedConfig) BranchesAndTypes(branches gitdomain.LocalBranchNames) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch)
	}
	return result
}

// ContainsLineage indicates whether this configuration contains any lineage entries.
func (self *UnvalidatedConfig) ContainsLineage() bool {
	return self.Lineage.Len() > 0
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *UnvalidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return branch == mainBranch
	}
	return false
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *UnvalidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.IsMainBranch(branch) || self.IsPerennialBranch(branch)
}

func (self *UnvalidatedConfig) IsOnline() bool {
	return self.Online().IsTrue()
}

func (self *UnvalidatedConfig) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	if slices.Contains(self.PerennialBranches, branch) {
		return true
	}
	if perennialRegex, has := self.PerennialRegex.Get(); has {
		return perennialRegex.MatchesBranch(branch)
	}
	return false
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, self.PerennialBranches...)
	}
	return self.PerennialBranches
}

func (self *UnvalidatedConfig) MatchesFeatureBranchRegex(branch gitdomain.LocalBranchName) bool {
	if featureRegex, has := self.FeatureRegex.Get(); has {
		return featureRegex.MatchesBranch(branch)
	}
	return false
}

func (self *UnvalidatedConfig) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *UnvalidatedConfig) Online() Online {
	return self.Offline.ToOnline()
}

func (self *UnvalidatedConfig) ShouldPushNewBranches() bool {
	return self.PushNewBranches.IsTrue()
}

// UnvalidatedBranchesAndTypes provides the types for the given branches.
// This method's name startes with "Unvalidated" to indicate that the types might be incomplete,
// and you should use ValidatedConfig.BranchesAndTypes if possible.
func (self *UnvalidatedConfig) UnvalidatedBranchesAndTypes(branches gitdomain.LocalBranchNames) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch)
	}
	return result
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() UnvalidatedConfig {
	return UnvalidatedConfig{
		Aliases:                  Aliases{},
		ContributionBranches:     gitdomain.NewLocalBranchNames(),
		CreatePrototypeBranches:  false,
		DefaultBranchType:        DefaultBranchType{BranchType: BranchTypeFeatureBranch},
		FeatureRegex:             None[FeatureRegex](),
		GitHubToken:              None[GitHubToken](),
		GitLabToken:              None[GitLabToken](),
		GitUserEmail:             None[GitUserEmail](),
		GitUserName:              None[GitUserName](),
		GiteaToken:               None[GiteaToken](),
		HostingOriginHostname:    None[HostingOriginHostname](),
		HostingPlatform:          None[HostingPlatform](),
		Lineage:                  NewLineage(),
		MainBranch:               None[gitdomain.LocalBranchName](),
		ObservedBranches:         gitdomain.NewLocalBranchNames(),
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

func NewUnvalidatedConfig(configFile Option[PartialConfig], globalGitConfig, localGitConfig PartialConfig) UnvalidatedConfig {
	var result PartialConfig
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		result = configFile
	}
	result = result.Merge(globalGitConfig)
	result = result.Merge(localGitConfig)
	return result.ToUnvalidatedConfig(DefaultConfig())
}
