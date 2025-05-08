package opcodes

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// MergeIntoCurrentBranch merges the branch with the given name into the current branch.
type MergeIntoCurrentBranch struct {
	BranchToMerge           gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeIntoCurrentBranch) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *MergeIntoCurrentBranch) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *MergeIntoCurrentBranch) Run(args shared.RunArgs) error {
	return args.Git.MergeBranchNoEdit(args.Frontend, self.BranchToMerge)
}
