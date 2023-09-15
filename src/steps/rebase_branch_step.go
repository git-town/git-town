package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	Branch      domain.BranchName
	previousSHA domain.SHA
	currentSHA  domain.SHA
	EmptyStep
}

func (step *RebaseBranchStep) CreateAbortSteps() []Step {
	return []Step{&AbortRebaseStep{}}
}

func (step *RebaseBranchStep) CreateContinueSteps() []Step {
	return []Step{&ContinueRebaseStep{}}
}

func (step *RebaseBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetCurrentBranchToSHAStep{Hard: true, MustHaveSHA: step.currentSHA, SetToSHA: step.previousSHA}}, nil
}

func (step *RebaseBranchStep) Run(args RunArgs) error {
	var err error
	step.previousSHA, err = args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	err = args.Runner.Frontend.Rebase(step.Branch)
	if err != nil {
		return err
	}
	step.currentSHA, err = args.Runner.Backend.CurrentSHA()
	return err
}
