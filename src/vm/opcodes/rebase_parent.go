package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RebaseParent rebases the given branch against the branch that is its parent at runtime.
type RebaseParent struct {
	CurrentBranch               gitdomain.LocalBranchName
	ParentActiveInOtherWorktree bool
	undeclaredOpcodeMethods     `exhaustruct:"optional"`
}

func (self *RebaseParent) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (self *RebaseParent) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueRebase{},
	}
}

func (self *RebaseParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.CurrentBranch).Get()
	if !hasParent {
		return nil
	}
	var branchToRebase gitdomain.BranchName
	if self.ParentActiveInOtherWorktree {
		branchToRebase = parent.TrackingBranch().BranchName()
	} else {
		branchToRebase = parent.BranchName()
	}
	return args.Frontend.Rebase(branchToRebase)
}
