package config

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v10/src/messages"
)

// PullBranchStrategy defines legal values for the "pull-branch-strategy" configuration setting.
type PullBranchStrategy struct {
	name string
}

func (self PullBranchStrategy) String() string { return self.name }

var (
	PullBranchStrategyMerge  = PullBranchStrategy{"merge"}  //nolint:gochecknoglobals
	PullBranchStrategyRebase = PullBranchStrategy{"rebase"} //nolint:gochecknoglobals
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
