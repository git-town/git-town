// Package steps defines the individual CLI operations that Git Town can execute.
// All steps implement the Step interface defined in step.go.
// Git Town doesn't execute steps directly.
// It organizes all Step instances it wants to perform in a StepList and executes that StepList.
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
	CreateUndoStep(*git.InternalCommands) (Step, error)

	// CreateAutomaticAbortError provides the error message to display when this step
	// cause the command to automatically abort.
	CreateAutomaticAbortError() error

	// Run executes this step.
	Run(repo *git.ProdRepo, connector hosting.Connector) error

	// ShouldAutomaticallyAbortOnError indicates whether this step should
	// cause the command to automatically abort if it errors.
	// When true, automatically runs the abort logic and leaves the user where they started.
	// When false, stops execution to let the user fix the issue and continue or manually abort.
	ShouldAutomaticallyAbortOnError() bool
}
