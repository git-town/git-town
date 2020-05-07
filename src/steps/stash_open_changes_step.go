package steps

import (
	"github.com/git-town/git-town/src/script"
)

// StashOpenChangesStep stores all uncommitted changes on the Git stash.
type StashOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStep returns the undo step for this step.
func (step *StashOpenChangesStep) CreateUndoStep() Step {
	return &RestoreOpenChangesStep{}
}

// Run executes this step.
func (step *StashOpenChangesStep) Run() error {
	err := script.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return script.RunCommand("git", "stash")
}
