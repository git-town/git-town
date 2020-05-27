package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStep returns the undo step for this step.
func (step *RestoreOpenChangesStep) CreateUndoStep() Step {
	return &StashOpenChangesStep{}
}

// Run executes this step.
func (step *RestoreOpenChangesStep) Run(repo *git.ProdRepo) error {
	return script.RunCommand("git", "stash", "pop")
}
