package steps

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type CommitOpenChangesStep struct{}

func (step CommitOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step CommitOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CommitOpenChangesStep) CreateUndoStep() Step {
	branchName := git.GetCurrentBranchName()
	return ResetToShaStep{Sha: git.GetBranchSha(branchName)}
}

func (step CommitOpenChangesStep) Run() error {
	err := script.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return script.RunCommand("git", "commit", "-m", fmt.Sprintf("'WIP on %s'", git.GetCurrentBranchName()))
}
