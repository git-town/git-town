package configdomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
)

// BranchTypes answers questions about whether branches are long-lived or not.
type BranchTypes struct {
	MainBranch        gitdomain.LocalBranchName
	PerennialBranches gitdomain.LocalBranchNames
}

func (self BranchTypes) IsFeatureBranch(branch gitdomain.LocalBranchName) bool {
	return !self.IsMainBranch(branch) && !self.IsPerennialBranch(branch)
}

func (self BranchTypes) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}

func (self BranchTypes) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.PerennialBranches, branch)
}

func (self BranchTypes) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.MainBranch}, self.PerennialBranches...)
}

func EmptyBranchTypes() BranchTypes {
	return BranchTypes{
		MainBranch:        gitdomain.EmptyLocalBranchName(),
		PerennialBranches: gitdomain.LocalBranchNames{},
	}
}
