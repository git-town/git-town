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
	previousSHA domain.SHA
	EmptyStep
}

func (step *CommitOpenChangesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetToSHAStep{SHA: step.previousSHA, Hard: false}}, nil
}

func (step *CommitOpenChangesStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.previousSHA, err = run.Backend.CurrentSHA()
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
