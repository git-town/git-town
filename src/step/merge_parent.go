package step

import "github.com/git-town/git-town/v9/src/domain"

// MergeParent merges the current parent of the current branch into the current branch.
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
	if parent == step.CurrentBranch {
		return nil
	}
	err := args.Runner.Frontend.MergeBranchNoEdit(parent.BranchName())
	return err
}
