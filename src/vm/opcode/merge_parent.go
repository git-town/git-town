package opcode

import "github.com/git-town/git-town/v9/src/domain"

// MergeParent merges the branch that at runtime is the parent branch of the given branch into the given branch.
type MergeParent struct {
	CurrentBranch domain.LocalBranchName
	BaseOpcode
}

func (step *MergeParent) CreateAbortProgram() []Opcode {
	return []Opcode{
		&AbortMerge{},
	}
}

func (step *MergeParent) CreateContinueProgram() []Opcode {
	return []Opcode{
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
