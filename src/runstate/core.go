// Package runstate represents the current execution status of a Git Town command,
// i.e. which steps to execute via the `runvm`.
package runstate

import (
	"time"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/steps"
)

// UnfinishedRunStateDetails has details about an unfinished run state.
type UnfinishedRunStateDetails struct {
	CanSkip   bool
	EndBranch domain.LocalBranchName
	EndTime   time.Time
}

func isCheckoutStep(step steps.Step) bool {
	return gohacks.TypeName(step) == "CheckoutStep"
}
