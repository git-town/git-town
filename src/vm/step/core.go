// Package steps defines the individual operations that the Git Town VM can execute.
// All steps implement the Step interface defined in step.go.
package step

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// Step represents a dedicated CLI activity.
// Git Town commands consist of many Step instances.
// Steps implement the command pattern (https://en.wikipedia.org/wiki/Command_pattern)
// and can provide Steps to continue, abort, and undo them.
type Step interface {
	// CreateAbortProgram provides the abort step for this step.
	CreateAbortProgram() []Step

	// CreateContinueProgram provides the continue step for this step.
	CreateContinueProgram() []Step

	// CreateAutomaticAbortError provides the error message to display when this step
	// cause the command to automatically abort.
	CreateAutomaticAbortError() error

	// Run executes this step.
	Run(args RunArgs) error

	// ShouldAutomaticallyAbortOnError indicates whether this step should
	// cause the command to automatically abort if it errors.
	// When true, automatically runs the abort logic and leaves the user where they started.
	// When false, stops execution to let the user fix the issue and continue or manually abort.
	ShouldAutomaticallyAbortOnError() bool
}

type RunArgs struct {
	AddSteps                        func(...Step) // AddSteps allows currently executing steps to prepend additional steps onto the currently executing step list.
	Connector                       hosting.Connector
	Lineage                         config.Lineage
	RegisterUndoablePerennialCommit func(domain.SHA)
	Runner                          *git.ProdRunner
	UpdateInitialBranchLocalSHA     func(domain.LocalBranchName, domain.SHA) error
}
