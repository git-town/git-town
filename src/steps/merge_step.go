package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// MergeStep merges the branch with the given name into the current branch.
type MergeStep struct {
	EmptyStep
	Branch      string
	previousSha domain.SHA
}

func (step *MergeStep) CreateAbortStep() Step {
	return &AbortMergeStep{}
}

func (step *MergeStep) CreateContinueStep() Step {
	return &ContinueMergeStep{}
}

func (step *MergeStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetToShaStep{Hard: true, Sha: step.previousSha}}, nil
}

func (step *MergeStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.previousSha, err = run.Backend.CurrentSha()
	if err != nil {
		return err
	}
	return run.Frontend.MergeBranchNoEdit(step.Branch)
}
