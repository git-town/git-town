// Package runstate represents the current execution status of a Git Town command and allows to run or persist it to disk.
// This is used by the "abort", "continue", and "undo" commands.
// The central data structure is RunState.
package runstate

import (
	"time"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/steps"
)

// UnfinishedRunStateDetails has details about an unfinished run state.
type UnfinishedRunStateDetails struct {
	CanSkip   bool
	EndBranch domain.LocalBranchName
	EndTime   time.Time
}

func isCheckoutStep(step steps.Step) bool {
	return typeName(step) == "*CheckoutStep"
}
