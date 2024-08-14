package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// UnvalidatedConfig is the Git Town configuration as read from disk.
// It might be lacking essential information in case Git metadata and config files don't contain it.
// If you need this information, validate it into a ValidatedConfig.
type UnvalidatedConfig struct {
	Aliases                  Aliases
	ContributionBranches     gitdomain.LocalBranchNames
	CreatePrototypeBranches  CreatePrototypeBranches
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
	SyncFeatureStrategy      SyncFeatureStrategy
	SyncPerennialStrategy    SyncPerennialStrategy
	SyncPrototypeStrategy    SyncPrototypeStrategy
	SyncTags                 SyncTags
	SyncUpstream             SyncUpstream
}

func (self *UnvalidatedConfig) BranchType(branch gitdomain.LocalBranchName) BranchType {
	switch {
	case self.IsMainBranch(branch):
		return BranchTypeMainBranch
	case self.IsPerennialBranch(branch):
		return BranchTypePerennialBranch
	case self.IsContributionBranch(branch):
		return BranchTypeContributionBranch
	case self.IsObservedBranch(branch):
		return BranchTypeObservedBranch
	case self.IsParkedBranch(branch):
		return BranchTypeParkedBranch
	case self.IsPrototypeBranch(branch):
		return BranchTypePrototypeBranch
	}
	return BranchTypeFeatureBranch
}

// ContainsLineage indicates whether this configuration contains any lineage entries.
func (self *UnvalidatedConfig) ContainsLineage() bool {
	return self.Lineage.Len() > 0
}

func (self *UnvalidatedConfig) IsContributionBranch(branch gitdomain.LocalBranchName) bool {
	return slices.Contains(self.ContributionBranches, branch)
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

func (self *UnvalidatedConfig) IsObservedBranch(branch gitdomain.LocalBranchName) bool {
	return slices.Contains(self.ObservedBranches, branch)
}

func (self *UnvalidatedConfig) IsOnline() bool {
	return self.Online().IsTrue()
}

func (self *UnvalidatedConfig) IsParkedBranch(branch gitdomain.LocalBranchName) bool {
	return slices.Contains(self.ParkedBranches, branch)
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

func (self *UnvalidatedConfig) IsPrototypeBranch(branch gitdomain.LocalBranchName) bool {
	return slices.Contains(self.PrototypeBranches, branch)
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, self.PerennialBranches...)
	}
	return self.PerennialBranches
}

func (self *UnvalidatedConfig) MustKnowParent(branch gitdomain.LocalBranchName) bool {
	return !self.IsMainBranch(branch) && !self.IsPerennialBranch(branch) && !self.IsContributionBranch(branch) && !self.IsObservedBranch(branch)
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

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() UnvalidatedConfig {
	return UnvalidatedConfig{
		Aliases:                  Aliases{},
		ContributionBranches:     gitdomain.NewLocalBranchNames(),
		CreatePrototypeBranches:  false,
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
