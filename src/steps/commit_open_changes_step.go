package steps

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	EmptyStep
	previousSha string
}

func (step *CommitOpenChangesStep) CreateUndoStep(repo *git.PublicRepo) (Step, error) {
	return &ResetToShaStep{Sha: step.previousSha}, nil
}

func (step *CommitOpenChangesStep) Run(repo *git.PublicRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.CurrentSha()
	if err != nil {
		return err
	}
	err = repo.StageFiles("-A")
	if err != nil {
		return err
	}
	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		return err
	}
	return repo.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch))
}
