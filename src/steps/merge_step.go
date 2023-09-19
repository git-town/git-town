package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// MergeStep merges the branch with the given name into the current branch.
type MergeStep struct {
	Branch      domain.BranchName
	previousSHA domain.SHA
	EmptyStep
}

func (step *MergeStep) CreateAbortSteps() []Step {
	return []Step{&AbortMergeStep{}}
}

func (step *MergeStep) CreateContinueSteps() []Step {
	return []Step{&ContinueMergeStep{}}
}

func (step *MergeStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetCurrentBranchToSHAStep{Hard: true, SHA: step.previousSHA}}, nil
}

func (step *MergeStep) Run(args RunArgs) error {
	var err error
	step.previousSHA, err = args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	return args.Runner.Frontend.MergeBranchNoEdit(step.Branch)
}
