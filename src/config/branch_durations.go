package config

import (
	"github.com/git-town/git-town/v9/src/slice"
)

// BranchDurations answers questions about whether branches are long-lived or not.
type BranchDurations struct {
	MainBranch        string
	PerennialBranches []string
}

func (pb BranchDurations) IsFeatureBranch(branch string) bool {
	return branch != pb.MainBranch && !slice.Contains(pb.PerennialBranches, branch)
}

func (pb BranchDurations) IsMainBranch(branch string) bool {
	return branch == pb.MainBranch
}

func (pb BranchDurations) IsPerennialBranch(branch string) bool {
	return slice.Contains(pb.PerennialBranches, branch)
}

func EmptyBranchDurations() BranchDurations {
	return BranchDurations{
		MainBranch:        "",
		PerennialBranches: []string{},
	}
}
