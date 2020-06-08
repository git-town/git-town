package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
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
func (step *CommitOpenChangesStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	step.previousSha = git.GetCurrentSha()
	err := repo.Logging.StageFiles("-A")
	if err != nil {
		return err
	}
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	return repo.Logging.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch))
}
