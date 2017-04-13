package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

// CreateAbortStep returns the abort step for this step.
func (step DiscardOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step DiscardOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step DiscardOpenChangesStep) Run() error {
	return script.RunCommand("git", "reset", "--hard")
}
