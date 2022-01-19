package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// Step represents a dedicated activity within a Git Town command.
// Git Town commands are comprised of a number of steps that need to be executed.
type Step interface {
	CreateAbortStep() Step
	CreateContinueStep() Step

	// CreateUndoStep returns the undo step for this step.
	CreateUndoStep(*git.ProdRepo) (Step, error)
	CreateAutomaticAbortError() error

	// Run executes this step.
	Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error
	ShouldAutomaticallyAbortOnError() bool
}
