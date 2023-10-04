package domain

import "github.com/git-town/git-town/v9/src/gohacks/slice"

// BranchTypes answers questions about whether branches are long-lived or not.
type BranchTypes struct {
	MainBranch        LocalBranchName
	PerennialBranches LocalBranchNames
}

func (bts BranchTypes) IsFeatureBranch(branch LocalBranchName) bool {
	return !bts.IsMainBranch(branch) && !bts.IsPerennialBranch(branch)
}

func (bts BranchTypes) IsMainBranch(branch LocalBranchName) bool {
	return branch == bts.MainBranch
}

func (bts BranchTypes) IsPerennialBranch(branch LocalBranchName) bool {
	return slice.Contains(bts.PerennialBranches, branch)
}

func EmptyBranchTypes() BranchTypes {
	return BranchTypes{
		MainBranch:        LocalBranchName{},
		PerennialBranches: LocalBranchNames{},
	}
}
