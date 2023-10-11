package step

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
)

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
	fmt.Println("111111111111111111111 START MERGE PARENT")
	parent := args.Lineage.Parent(step.CurrentBranch)
	fmt.Println("111111111111111111111 CURRENT BRANCH", step.CurrentBranch)
	fmt.Println("111111111111111111111 PARENT", parent)
	if parent.IsEmpty() {
		return nil
	}
	if parent == step.CurrentBranch {
		return nil
	}
	err := args.Runner.Frontend.MergeBranchNoEdit(parent.BranchName())
	fmt.Println("111111111111111111111 FINISHED MERGE PARENT")
	return err
}
