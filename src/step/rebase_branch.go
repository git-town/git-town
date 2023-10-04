package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// RebaseBranch rebases the current branch
// against the branch with the given name.
type RebaseBranch struct {
	Branch domain.BranchName
	Empty
}

func (step *RebaseBranch) CreateAbortSteps() []Step {
	return []Step{&AbortRebase{}}
}

func (step *RebaseBranch) CreateContinueSteps() []Step {
	return []Step{&ContinueRebase{}}
}

func (step *RebaseBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.Rebase(step.Branch)
}
