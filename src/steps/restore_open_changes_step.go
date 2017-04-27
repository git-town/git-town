package steps

import (
	"github.com/Originate/git-town/src/script"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step RestoreOpenChangesStep) CreateUndoStepBeforeRun() Step {
	return StashOpenChangesStep{}
}

// Run executes this step.
func (step RestoreOpenChangesStep) Run() error {
	return script.RunCommand("git", "stash", "pop")
}
