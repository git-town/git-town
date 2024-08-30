package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// MergeParent merges the branch that at runtime is the parent branch of the given branch into the given branch.
type MergeParent struct {
	CurrentBranch               gitdomain.LocalBranchName
	ParentActiveInOtherWorktree bool
	undeclaredOpcodeMethods     `exhaustruct:"optional"`
}

func (self *MergeParent) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
	}
}

func (self *MergeParent) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueMerge{},
	}
}

func (self *MergeParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.CurrentBranch).Get()
	if !hasParent {
		return nil
	}
	var branchToMerge gitdomain.BranchName
	if self.ParentActiveInOtherWorktree {
		branchToMerge = parent.TrackingBranch().BranchName()
	} else {
		branchToMerge = parent.BranchName()
	}
	return args.Git.MergeBranchNoEdit(args.Frontend, branchToMerge)
}
