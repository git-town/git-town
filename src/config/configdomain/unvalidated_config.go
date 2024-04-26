package configdomain

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
)

// UnvalidatedConfig is the merged configuration to be used by Git Town commands.
type UnvalidatedConfig struct {
	Aliases                  Aliases
	ContributionBranches     gitdomain.LocalBranchNames
	GitHubToken              Option[GitHubToken]
	GitLabToken              Option[GitLabToken]
	GitUserEmail             GitUserEmail
	GitUserName              GitUserName
	GiteaToken               Option[GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	HostingPlatform          Option[HostingPlatform] // Some = override by user, None = auto-detect
	Lineage                  Lineage
	MainBranch               gitdomain.LocalBranchName
	ObservedBranches         gitdomain.LocalBranchNames
	Offline                  Offline
	ParkedBranches           gitdomain.LocalBranchNames
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PushHook                 PushHook
	PushNewBranches          PushNewBranches
	ShipDeleteTrackingBranch ShipDeleteTrackingBranch
	SyncBeforeShip           SyncBeforeShip
	SyncFeatureStrategy      SyncFeatureStrategy
	SyncPerennialStrategy    SyncPerennialStrategy
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
	}
	return BranchTypeFeatureBranch
}

// ContainsLineage indicates whether this configuration contains any lineage entries.
func (self *UnvalidatedConfig) ContainsLineage() bool {
	return len(self.Lineage) > 0
}

func (self *UnvalidatedConfig) IsContributionBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.ContributionBranches, branch)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *UnvalidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *UnvalidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.IsMainBranch(branch) || self.IsPerennialBranch(branch)
}

func (self *UnvalidatedConfig) IsObservedBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.ObservedBranches, branch)
}

func (self *UnvalidatedConfig) IsOnline() bool {
	return self.Online().Bool()
}

func (self *UnvalidatedConfig) IsParkedBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.ParkedBranches, branch)
}

func (self *UnvalidatedConfig) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	if slice.Contains(self.PerennialBranches, branch) {
		return true
	}
	if perennialRegex, has := self.PerennialRegex.Get(); has {
		return perennialRegex.MatchesBranch(branch)
	}
	return false
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.MainBranch}, self.PerennialBranches...)
}

// Merges the given PartialConfig into this configuration object.
func (self *UnvalidatedConfig) Merge(other PartialConfig) {
	for key, value := range other.Aliases {
		self.Aliases[key] = value
	}
	for child, parent := range other.Lineage {
		self.Lineage[child] = parent
	}
	self.ContributionBranches = append(self.ContributionBranches, other.ContributionBranches...)
	if other.HostingOriginHostname.IsSome() {
		self.HostingOriginHostname = other.HostingOriginHostname
	}
	if other.HostingPlatform.IsSome() {
		self.HostingPlatform = other.HostingPlatform
	}
	if other.GiteaToken.IsSome() {
		self.GiteaToken = other.GiteaToken
	}
	if other.GitHubToken.IsSome() {
		self.GitHubToken = other.GitHubToken
	}
	if other.GitLabToken.IsSome() {
		self.GitLabToken = other.GitLabToken
	}
	if email, has := other.GitUserEmail.Get(); has {
		self.GitUserEmail = email
	}
	if name, has := other.GitUserName.Get(); has {
		self.GitUserName = name
	}
	if branch, has := other.MainBranch.Get(); has {
		self.MainBranch = branch
	}
	if pushNewBranches, has := other.PushNewBranches.Get(); has {
		self.PushNewBranches = pushNewBranches
	}
	self.ObservedBranches = append(self.ObservedBranches, other.ObservedBranches...)
	if offline, has := other.Offline.Get(); has {
		self.Offline = offline
	}
	self.ParkedBranches = append(self.ParkedBranches, other.ParkedBranches...)
	self.PerennialBranches = append(self.PerennialBranches, other.PerennialBranches...)
	if other.PerennialRegex.IsSome() {
		self.PerennialRegex = other.PerennialRegex
	}
	if value, has := other.PushHook.Get(); has {
		self.PushHook = value
	}
	if value, has := other.ShipDeleteTrackingBranch.Get(); has {
		self.ShipDeleteTrackingBranch = value
	}
	if value, has := other.SyncBeforeShip.Get(); has {
		self.SyncBeforeShip = value
	}
	if value, has := other.SyncFeatureStrategy.Get(); has {
		self.SyncFeatureStrategy = value
	}
	if value, has := other.SyncPerennialStrategy.Get(); has {
		self.SyncPerennialStrategy = value
	}
	if other.SyncUpstream != nil {
		self.SyncUpstream = *other.SyncUpstream
	}
}

func (self *UnvalidatedConfig) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *UnvalidatedConfig) Online() Online {
	return self.Offline.ToOnline()
}

func (self *UnvalidatedConfig) ShouldPushNewBranches() bool {
	return self.PushNewBranches.Bool()
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() UnvalidatedConfig {
	return UnvalidatedConfig{
		Aliases:                  Aliases{},
		ContributionBranches:     gitdomain.NewLocalBranchNames(),
		GitHubToken:              None[GitHubToken](),
		GitLabToken:              None[GitLabToken](),
		GitUserEmail:             "",
		GitUserName:              "",
		GiteaToken:               None[GiteaToken](),
		HostingOriginHostname:    None[HostingOriginHostname](),
		HostingPlatform:          None[HostingPlatform](),
		Lineage:                  Lineage{},
		MainBranch:               gitdomain.EmptyLocalBranchName(),
		ObservedBranches:         gitdomain.NewLocalBranchNames(),
		Offline:                  false,
		ParkedBranches:           gitdomain.NewLocalBranchNames(),
		PerennialBranches:        gitdomain.NewLocalBranchNames(),
		PerennialRegex:           None[PerennialRegex](),
		PushHook:                 true,
		PushNewBranches:          false,
		ShipDeleteTrackingBranch: true,
		SyncBeforeShip:           false,
		SyncFeatureStrategy:      SyncFeatureStrategyMerge,
		SyncPerennialStrategy:    SyncPerennialStrategyRebase,
		SyncUpstream:             true,
	}
}

func NewFullConfig(configFile *PartialConfig, globalGitConfig, localGitConfig PartialConfig) UnvalidatedConfig {
	result := DefaultConfig()
	if configFile != nil {
		result.Merge(*configFile)
	}
	result.Merge(globalGitConfig)
	result.Merge(localGitConfig)
	return result
}
