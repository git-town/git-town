package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	previousSHA domain.SHA
	EmptyStep
}

func (step *CommitOpenChangesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetCurrentBranchToSHAStep{SHA: step.previousSHA, Hard: false}}, nil
}

func (step *CommitOpenChangesStep) Run(args RunArgs) error {
	var err error
	step.previousSHA, err = args.Run.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	err = args.Run.Frontend.StageFiles("-A")
	if err != nil {
		return err
	}
	currentBranch, err := args.Run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	return args.Run.Frontend.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch))
}
