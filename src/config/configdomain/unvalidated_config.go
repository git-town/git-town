package configdomain

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
)

// UnvalidatedConfig is the Git Town configuration as read from disk.
// It might be lacking essential information in case Git metadata and config files don't contain it.
// If you need this information, validate it into a ValidatedConfig.
type UnvalidatedConfig struct {
	Aliases                  Aliases
	ContributionBranches     gitdomain.LocalBranchNames
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
	return slice.Contains(self.ContributionBranches, branch)
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

func (self *UnvalidatedConfig) IsPrototypeBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.PrototypeBranches, branch)
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, self.PerennialBranches...)
	}
	return self.PerennialBranches
}

// Merges the given PartialConfig into this configuration object.
func (self *UnvalidatedConfig) Merge(other PartialConfig) {
	for key, value := range other.Aliases {
		self.Aliases[key] = value
	}
	for _, entry := range other.Lineage.Entries() {
		self.Lineage.Add(entry.Child, entry.Parent)
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
	if other.GitUserEmail.IsSome() {
		self.GitUserEmail = other.GitUserEmail
	}
	if other.GitUserName.IsSome() {
		self.GitUserName = other.GitUserName
	}
	if branch, has := other.MainBranch.Get(); has {
		self.MainBranch = Some(branch)
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
	self.PrototypeBranches = append(self.PrototypeBranches, other.PrototypeBranches...)
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
	if value, has := other.SyncUpstream.Get(); has {
		self.SyncUpstream = value
	}
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
	return self.PushNewBranches.Bool()
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() UnvalidatedConfig {
	return UnvalidatedConfig{
		Aliases:                  Aliases{},
		ContributionBranches:     gitdomain.NewLocalBranchNames(),
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
		SyncBeforeShip:           false,
		SyncFeatureStrategy:      SyncFeatureStrategyMerge,
		SyncPerennialStrategy:    SyncPerennialStrategyRebase,
		SyncUpstream:             true,
	}
}

func NewUnvalidatedConfig(configFile Option[PartialConfig], globalGitConfig, localGitConfig PartialConfig) UnvalidatedConfig {
	result := DefaultConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		result.Merge(configFile)
	}
	result.Merge(globalGitConfig)
	result.Merge(localGitConfig)
	return result
}
