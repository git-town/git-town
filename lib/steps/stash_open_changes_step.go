package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// StashOpenChangesStep stores all uncommitted changes on the Git stash.
type StashOpenChangesStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step StashOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step StashOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step StashOpenChangesStep) CreateUndoStep() Step {
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
