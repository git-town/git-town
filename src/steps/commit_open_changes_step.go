package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	beforeSHA domain.SHA
	afterSHA  domain.SHA
	EmptyStep
}

func (step *CommitOpenChangesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetCurrentBranchToSHAStep{MustHaveSHA: step.afterSHA, SetToSHA: step.beforeSHA, Hard: false}}, nil
}

func (step *CommitOpenChangesStep) Run(args RunArgs) error {
	var err error
	step.beforeSHA, err = args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	err = args.Runner.Frontend.StageFiles("-A")
	if err != nil {
		return err
	}
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	err = args.Runner.Frontend.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch))
	if err != nil {
		return err
	}
	step.afterSHA, err = args.Runner.Backend.CurrentSHA()
	return err
}
