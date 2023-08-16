package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	EmptyStep
	previousSha domain.SHA
}

func (step *CommitOpenChangesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetToShaStep{Sha: step.previousSha}}, nil
}

func (step *CommitOpenChangesStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.previousSha, err = run.Backend.CurrentSha()
	if err != nil {
		return err
	}
	err = run.Frontend.StageFiles("-A")
	if err != nil {
		return err
	}
	currentBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	return run.Frontend.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch))
}
