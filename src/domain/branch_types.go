package domain

import "github.com/git-town/git-town/v11/src/gohacks/slice"

// BranchTypes answers questions about whether branches are long-lived or not.
type BranchTypes struct {
	MainBranch        LocalBranchName
	PerennialBranches LocalBranchNames
}

func (self BranchTypes) IsFeatureBranch(branch LocalBranchName) bool {
	return !self.IsMainBranch(branch) && !self.IsPerennialBranch(branch)
}

func (self BranchTypes) IsMainBranch(branch LocalBranchName) bool {
	return branch == self.MainBranch
}

func (self BranchTypes) IsPerennialBranch(branch LocalBranchName) bool {
	return slice.Contains(self.PerennialBranches, branch)
}

func (self BranchTypes) MainAndPerennials() LocalBranchNames {
	return append(LocalBranchNames{self.MainBranch}, self.PerennialBranches...)
}

func EmptyBranchTypes() BranchTypes {
	return BranchTypes{
		MainBranch:        EmptyLocalBranchName(),
		PerennialBranches: LocalBranchNames{},
	}
}
