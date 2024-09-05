package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
)

// ValidatedConfig is Git Town configuration where all essential values are guaranteed to exist and have meaningful values.
// This is ensured by querying from the user if needed.
type ValidatedConfig struct {
	GitUserEmail GitUserEmail
	GitUserName  GitUserName
	MainBranch   gitdomain.LocalBranchName
	*UnvalidatedConfig
}

func (self *ValidatedConfig) BranchType(branch gitdomain.LocalBranchName) BranchType {
	if branch == self.MainBranch {
		return BranchTypeMainBranch
	}
	if slices.Contains(self.PerennialBranches, branch) {
		return BranchTypePerennialBranch
	}
	if perennialRegex, hasPerennialRegex := self.PerennialRegex.Get(); hasPerennialRegex {
		if perennialRegex.MatchesBranch(branch) {
			return BranchTypePerennialBranch
		}
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
	return BranchTypeFeatureBranch
}

func (self *ValidatedConfig) BranchesAndTypes(branches gitdomain.LocalBranchNames) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch)
	}
	return result
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

func (self *ValidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.MainBranch}, self.PerennialBranches...)
}

// provides this collection without the perennial branch at the root
func (self ValidatedConfig) RemovePerennials(stack gitdomain.LocalBranchNames) gitdomain.LocalBranchNames {
	if len(stack) == 0 {
		return stack
	}
	result := make(gitdomain.LocalBranchNames, 0, len(stack)-1)
	for _, branch := range stack {
		if !self.IsMainOrPerennialBranch(branch) {
			result = append(result, branch)
		}
	}
	return result
}
