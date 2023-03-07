package config

import "fmt"

// PullBranchStrategy defines legal values for the "pull-branch-strategy" configuration setting.
type PullBranchStrategy string

const (
	PullBranchStrategyMerge  PullBranchStrategy = "merge"
	PullBranchStrategyRebase PullBranchStrategy = "rebase"
)

func NewPullBranchStrategy(text string) (PullBranchStrategy, error) {
	switch text {
	case "merge":
		return PullBranchStrategyMerge, nil
	case "rebase", "":
		return PullBranchStrategyRebase, nil
	default:
		return PullBranchStrategyMerge, fmt.Errorf("unknown pull branch strategy: %q", text)
	}
}

func (pbs PullBranchStrategy) String() string {
	return string(pbs)
}
