package steps

import (
	"github.com/Originate/git-town/src/script"
)

// StashOpenChangesStep stores all uncommitted changes on the Git stash.
type StashOpenChangesStep struct {
	NoOpStep
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *StashOpenChangesStep) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&RestoreOpenChangesStep{})
}

// Run executes this step.
func (step *StashOpenChangesStep) Run() error {
	err := script.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return script.RunCommand("git", "stash")
}
