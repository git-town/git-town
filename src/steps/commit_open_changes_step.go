package steps

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	NoOpStep

	previousSha string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *CommitOpenChangesStep) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&ResetToShaStep{Sha: step.previousSha})
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
