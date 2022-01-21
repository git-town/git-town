package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// Step represents a dedicated CLI activity.
// Git Town commands consist of many Step instances.
// Steps implement the command pattern (https://en.wikipedia.org/wiki/Command_pattern)
// and can provide Steps to continue, abort, and undo them.
type Step interface {
	// CreateAbortStep provides the abort step for this step.
	CreateAbortStep() Step

	// CreateContinueStep provides the continue step for this step.
	CreateContinueStep() Step

	// CreateUndoStep provides the undo step for this step.
	CreateUndoStep(*git.ProdRepo) (Step, error)

	// CreateAutomaticAbortError provides the error message to display when this step
	// cause the command to automatically abort.
	CreateAutomaticAbortError() error

	// Run executes this step.
	Run(repo *git.ProdRepo, driver hosting.Driver) error

	// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
	// automatically abort if it errors.
	ShouldAutomaticallyAbortOnError() bool
}
