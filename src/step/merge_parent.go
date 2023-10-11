package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// MergeParent merges the branch that at runtime is the parent branch into the current branch.
type MergeParent struct {
	CurrentBranch domain.LocalBranchName
	Empty
}

func (step *MergeParent) CreateAbortSteps() []Step {
	return []Step{
		&AbortMerge{},
	}
}

func (step *MergeParent) CreateContinueSteps() []Step {
	return []Step{
		&ContinueMerge{},
	}
}

func (step *MergeParent) Run(args RunArgs) error {
	parent := args.Lineage.Parent(step.CurrentBranch)
	if parent.IsEmpty() {
		return nil
	}
	return args.Runner.Frontend.MergeBranchNoEdit(parent.BranchName())
}
