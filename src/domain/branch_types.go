package domain

import "github.com/git-town/git-town/v9/src/slice"

// BranchTypes answers questions about whether branches are long-lived or not.
type BranchTypes struct {
	MainBranch        LocalBranchName
	PerennialBranches LocalBranchNames
}

func (pb BranchTypes) IsFeatureBranch(branch LocalBranchName) bool {
	return !pb.IsMainBranch(branch) && !pb.IsPerennialBranch(branch)
}

func (pb BranchTypes) IsMainBranch(branch LocalBranchName) bool {
	return branch == pb.MainBranch
}

func (pb BranchTypes) IsPerennialBranch(branch LocalBranchName) bool {
	return slice.Contains(pb.PerennialBranches, branch)
}

func EmptyBranchDurations() BranchTypes {
	return BranchTypes{
		MainBranch:        LocalBranchName{},
		PerennialBranches: LocalBranchNames{},
	}
}
