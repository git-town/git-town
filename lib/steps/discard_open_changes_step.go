package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	NoExpectedError
	NoUndoStep
}

// Run executes this step.
func (step DiscardOpenChangesStep) Run() error {
	return script.RunCommand("git", "reset", "--hard")
}
