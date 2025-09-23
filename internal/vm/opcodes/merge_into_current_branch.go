package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// MergeIntoCurrentBranch merges the branch with the given name into the current branch.
type MergeIntoCurrentBranch struct {
	BranchToMerge gitdomain.BranchName
}

func (self *MergeIntoCurrentBranch) Abort() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *MergeIntoCurrentBranch) Continue() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *MergeIntoCurrentBranch) Run(args shared.RunArgs) error {
	return args.Git.MergeBranchNoEdit(args.Frontend, self.BranchToMerge)
}
