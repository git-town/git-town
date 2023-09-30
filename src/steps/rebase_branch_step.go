package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	Branch domain.BranchName
	EmptyStep
}

func (step *RebaseBranchStep) CreateAbortSteps() []Step {
	return []Step{&AbortRebaseStep{}}
}

func (step *RebaseBranchStep) CreateContinueSteps() []Step {
	return []Step{&ContinueRebaseStep{}}
}

func (step *RebaseBranchStep) Run(args RunArgs) error {
	return args.Runner.Frontend.Rebase(step.Branch)
}
