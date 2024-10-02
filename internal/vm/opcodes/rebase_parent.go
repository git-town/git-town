package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// RebaseParent rebases the given branch against the branch that is its parent at runtime.
type RebaseParent struct {
	Parent                      gitdomain.BranchName
	ParentActiveInOtherWorktree bool
	undeclaredOpcodeMethods     `exhaustruct:"optional"`
}

func (self *RebaseParent) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortRebase{},
	}
}

func (self *RebaseParent) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueRebaseIfNeeded{},
	}
}

func (self *RebaseParent) Run(args shared.RunArgs) error {
	return args.Git.Rebase(args.Frontend, self.Parent)
}
