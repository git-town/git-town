package config

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/slice"
)

// BranchDurations answers questions about whether branches are long-lived or not.
type BranchDurations struct {
	MainBranch        domain.LocalBranchName
	PerennialBranches domain.LocalBranchNames
}

func (pb BranchDurations) IsFeatureBranch(branch domain.LocalBranchName) bool {
	return branch != pb.MainBranch && !slice.Contains(pb.PerennialBranches, branch)
}

func (pb BranchDurations) IsMainBranch(branch domain.LocalBranchName) bool {
	return branch == pb.MainBranch
}

func (pb BranchDurations) IsPerennialBranch(branch domain.LocalBranchName) bool {
	return slice.Contains(pb.PerennialBranches, branch)
}

func EmptyBranchDurations() BranchDurations {
	return BranchDurations{
		MainBranch:        domain.LocalBranchName{},
		PerennialBranches: domain.LocalBranchNames{},
	}
}
