package opcodes

import (
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// RebaseParent rebases the given branch against the branch that is its parent at runtime.
type RebaseParent struct {
	CurrentBranch               gitdomain.LocalBranchName
	ParentActiveInOtherWorktree bool
	undeclaredOpcodeMethods
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
	parent := args.Lineage.Parent(self.CurrentBranch)
	if parent.IsEmpty() {
		return nil
	}
	var branchToRebase gitdomain.BranchName
	if self.ParentActiveInOtherWorktree {
		branchToRebase = parent.TrackingBranch().BranchName()
	} else {
		branchToRebase = parent.BranchName()
	}
	return args.Runner.Frontend.Rebase(branchToRebase)
}
