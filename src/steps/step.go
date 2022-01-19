package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// Step represents a dedicated activity within a Git Town command.
// Git Town commands are comprised of a number of steps that need to be executed.
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
	Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error

	// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
	// automatically abort if it errors.
	ShouldAutomaticallyAbortOnError() bool
}
