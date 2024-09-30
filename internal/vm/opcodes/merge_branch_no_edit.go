package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// MergeParent merges the branch that at runtime is the parent branch of the given branch into the given branch.
type MergeBranchNoEdit struct {
	Branch                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeBranchNoEdit) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
	}
}

func (self *MergeBranchNoEdit) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueMerge{},
	}
}

func (self *MergeBranchNoEdit) Run(args shared.RunArgs) error {
	return args.Git.MergeBranchNoEdit(args.Frontend, self.Branch)
}
