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

func (step *CommitOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &ResetToShaStep{Sha: step.previousSha}, nil
}

func (step *CommitOpenChangesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	err = repo.Logging.StageFiles("-A")
	if err != nil {
		return err
	}
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	return repo.Logging.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch))
}
