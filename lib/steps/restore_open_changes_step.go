package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step RestoreOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step RestoreOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step RestoreOpenChangesStep) CreateUndoStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step RestoreOpenChangesStep) Run() error {
	return script.RunCommand("git", "stash", "pop")
}
