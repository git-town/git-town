package steps

import (
	"fmt"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	EmptyStep
}

func (step *CommitOpenChangesStep) Run(args RunArgs) error {
	err := args.Runner.Frontend.StageFiles("-A")
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
