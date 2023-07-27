package config

import "github.com/git-town/git-town/v9/src/stringslice"

type BranchDurations struct {
	MainBranch        string
	PerennialBranches []string
}

func (pb BranchDurations) IsFeatureBranch(branch string) bool {
	return branch != pb.MainBranch && !stringslice.Contains(pb.PerennialBranches, branch)
}

func (pb BranchDurations) IsMainBranch(branch string) bool {
	return branch == pb.MainBranch
}

func (pb BranchDurations) IsPerennialBranch(branch string) bool {
	return stringslice.Contains(pb.PerennialBranches, branch)
}
