package steps

import (
	"github.com/Originate/git-town/src/script"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *RestoreOpenChangesStep) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&StashOpenChangesStep{})
}

// Run executes this step.
func (step *RestoreOpenChangesStep) Run() error {
	return script.RunCommand("git", "stash", "pop")
}
