package domain

import (
	"github.com/git-town/git-town/v9/src/slice"
)

// BranchDurations answers questions about whether branches are long-lived or not.
// TODO: rename to Perennials
type BranchDurations struct {
	MainBranch        LocalBranchName
	PerennialBranches LocalBranchNames
}

func (pb BranchDurations) IsFeatureBranch(branch LocalBranchName) bool {
	return !pb.IsMainBranch(branch) && !pb.IsPerennialBranch(branch)
}

func (pb BranchDurations) IsMainBranch(branch LocalBranchName) bool {
	return branch == pb.MainBranch
}

func (pb BranchDurations) IsPerennialBranch(branch LocalBranchName) bool {
	return slice.Contains(pb.PerennialBranches, branch)
}

func EmptyBranchDurations() BranchDurations {
	return BranchDurations{
		MainBranch:        LocalBranchName{},
		PerennialBranches: LocalBranchNames{},
	}
}
