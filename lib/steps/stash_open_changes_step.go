package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// StashOpenChangesStep stores all uncommitted changes on the Git stash.
type StashOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step StashOpenChangesStep) CreateUndoStepBeforeRun() Step {
	return RestoreOpenChangesStep{}
}

// Run executes this step.
func (step StashOpenChangesStep) Run() error {
	err := script.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return script.RunCommand("git", "stash")
}
