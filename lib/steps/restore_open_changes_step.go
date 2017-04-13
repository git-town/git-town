package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

// CreateAbortStep returns the abort step for this step.
func (step RestoreOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step RestoreOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step RestoreOpenChangesStep) Run() error {
	return script.RunCommand("git", "stash", "pop")
}
