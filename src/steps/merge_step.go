package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// MergeStep merges the branch with the given name into the current branch.
type MergeStep struct {
	Branch domain.BranchName
	EmptyStep
}

func (step *MergeStep) CreateAbortSteps() []Step {
	return []Step{&AbortMergeStep{}}
}

func (step *MergeStep) CreateContinueSteps() []Step {
	return []Step{&ContinueMergeStep{}}
}

func (step *MergeStep) Run(args RunArgs) error {
	return args.Runner.Frontend.MergeBranchNoEdit(step.Branch)
}
