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
	EmptyStep
}

func (step *RebaseBranchStep) CreateAbortSteps() []Step {
	return []Step{&AbortRebaseStep{}}
}

func (step *RebaseBranchStep) CreateContinueSteps() []Step {
	return []Step{&ContinueRebaseStep{}}
}

func (step *RebaseBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetCurrentBranchToSHAStep{Hard: true, SHA: step.previousSHA}}, nil
}

func (step *RebaseBranchStep) Run(args RunArgs) error {
	var err error
	step.previousSHA, err = args.Run.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	return args.Run.Frontend.Rebase(step.Branch)
}
