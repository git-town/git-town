package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// RebaseParent rebases the current branch
// against the branch with the given name.
type RebaseParent struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *RebaseParent) CreateAbortSteps() []Step {
	return []Step{&AbortRebase{}}
}

func (step *RebaseParent) CreateContinueSteps() []Step {
	return []Step{&ContinueRebase{}}
}

func (step *RebaseParent) Run(args RunArgs) error {
	parent := args.Lineage.Parent(step.Branch)
	return args.Runner.Frontend.Rebase(parent.BranchName())
}
