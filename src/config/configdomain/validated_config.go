package configdomain

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/validate"
)

// ValidatedConfig is validated UnvalidatedConfig
type ValidatedConfig struct {
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

func NewValidatedConfig(unvalidated UnvalidatedConfig) ValidatedConfig {
	validatedMainBranch, validatedPerennialBranches, err := validate.MainAndPerennials(unvalidated.MainBranch, unvalidated.PerennialBranches)

	validatedGitUserEmail := validateGitUserEmail(unvalidated.GitUserEmail)
	validatedGitUserName := validateGitUserName(unvalidated.GitUserName)
	validatedLineage := validateLineage(unvalidated.Lineage)
	return ValidatedConfig{
		Aliases:                  unvalidated.Aliases,
		ContributionBranches:     unvalidated.ContributionBranches,
		GitHubToken:              unvalidated.GitHubToken,
		GitLabToken:              unvalidated.GitLabToken,
		GitUserEmail:             validatedGitUserEmail,
		GitUserName:              validatedGitUserName,
		GiteaToken:               unvalidated.GiteaToken,
		HostingOriginHostname:    unvalidated.HostingOriginHostname,
		HostingPlatform:          unvalidated.HostingPlatform,
		Lineage:                  validatedLineage,
		MainBranch:               validatedMainBranch,
		ObservedBranches:         unvalidated.ObservedBranches,
		Offline:                  unvalidated.Offline,
		ParkedBranches:           unvalidated.ParkedBranches,
		PerennialBranches:        validatedPerennialBranches,
		PerennialRegex:           unvalidated.PerennialRegex,
		PushHook:                 unvalidated.PushHook,
		PushNewBranches:          unvalidated.PushNewBranches,
		ShipDeleteTrackingBranch: unvalidated.ShipDeleteTrackingBranch,
		SyncBeforeShip:           unvalidated.SyncBeforeShip,
		SyncFeatureStrategy:      unvalidated.SyncFeatureStrategy,
		SyncPerennialStrategy:    unvalidated.SyncPerennialStrategy,
		SyncUpstream:             unvalidated.SyncUpstream,
	}
}

func (self *ValidatedConfig) BranchType(branch gitdomain.LocalBranchName) BranchType {
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
func (self *ValidatedConfig) ContainsLineage() bool {
	return len(self.Lineage) > 0
}

func (self *ValidatedConfig) IsContributionBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.ContributionBranches, branch)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *ValidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *ValidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.IsMainBranch(branch) || self.IsPerennialBranch(branch)
}

func (self *ValidatedConfig) IsObservedBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.ObservedBranches, branch)
}

func (self *ValidatedConfig) IsOnline() bool {
	return self.Online().Bool()
}

func (self *ValidatedConfig) IsParkedBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.ParkedBranches, branch)
}

func (self *ValidatedConfig) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	if slice.Contains(self.PerennialBranches, branch) {
		return true
	}
	if perennialRegex, hasPerennialRegex := self.PerennialRegex.Get(); hasPerennialRegex {
		return perennialRegex.MatchesBranch(branch)
	}
	return false
}

func (self *ValidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.MainBranch}, self.PerennialBranches...)
}

func (self *ValidatedConfig) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *ValidatedConfig) Online() Online {
	return self.Offline.ToOnline()
}

func (self *ValidatedConfig) ShouldPushNewBranches() bool {
	return self.PushNewBranches.Bool()
}
