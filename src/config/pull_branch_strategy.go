package config

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/messages"
)

// PullBranchStrategy defines legal values for the "pull-branch-strategy" configuration setting.
type PullBranchStrategy struct {
	name string
}

func (p PullBranchStrategy) String() string { return p.name }

var (
	PullBranchStrategyMerge  = PullBranchStrategy{"merge"}
	PullBranchStrategyRebase = PullBranchStrategy{"rebase"}
)

func NewPullBranchStrategy(text string) (PullBranchStrategy, error) {
	switch strings.ToLower(text) {
	case "merge":
		return PullBranchStrategyMerge, nil
	case "rebase", "":
		return PullBranchStrategyRebase, nil
	default:
		return PullBranchStrategyMerge, fmt.Errorf(messages.ConfigPullbranchStrategyUnknown, text)
	}
}
