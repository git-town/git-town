package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// MergeStep merges the branch with the given name into the current branch.
type MergeStep struct {
	Branch    domain.BranchName
	beforeSHA domain.SHA
	afterSHA  domain.SHA
	EmptyStep
}

func (step *MergeStep) CreateAbortSteps() []Step {
	return []Step{&AbortMergeStep{}}
}

func (step *MergeStep) CreateContinueSteps() []Step {
	return []Step{&ContinueMergeStep{}}
}

func (step *MergeStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetCurrentBranchToSHAStep{MustHaveSHA: step.afterSHA, SetToSHA: step.beforeSHA, Hard: true}}, nil
}

func (step *MergeStep) Run(args RunArgs) error {
	var err error
	step.beforeSHA, err = args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	err = args.Runner.Frontend.MergeBranchNoEdit(step.Branch)
	if err != nil {
		return err
	}
	step.afterSHA, err = args.Runner.Backend.CurrentSHA()
	return err
}
