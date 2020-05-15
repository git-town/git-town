package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	NoOpStep

	previousSha string
}

// CreateUndoStep returns the undo step for this step.
func (step *CommitOpenChangesStep) CreateUndoStep() Step {
	return &ResetToShaStep{Sha: step.previousSha}
}

// Run executes this step.
func (step *CommitOpenChangesStep) Run() error {
	step.previousSha = git.GetCurrentSha()
	err := script.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return script.RunCommand("git", "commit", "-m", fmt.Sprintf("WIP on %s", git.GetCurrentBranchName()))
}
