package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	EmptyStep
	Branch      string
	previousSha string
}

func (step *RebaseBranchStep) CreateAbortStep() Step {
	return &AbortRebaseStep{}
}

func (step *RebaseBranchStep) CreateContinueStep() Step {
	return &ContinueRebaseStep{}
}

func (step *RebaseBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetToShaStep{Hard: true, Sha: step.previousSha}}, nil
}

func (step *RebaseBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.previousSha, err = run.Backend.CurrentSha()
	if err != nil {
		return err
	}
	return run.Frontend.Rebase(step.Branch)
}
