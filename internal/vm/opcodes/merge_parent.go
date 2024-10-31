package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// MergeParent merges the given parent branch into the current branch.
type MergeParent struct {
	Parent                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParent) Run(args shared.RunArgs) error {
	err := args.Git.MergeBranchNoEdit(args.Frontend, self.Parent)
	if err != nil {
		args.PrependOpcodes(&ConflictPhantomDetect{
			ParentBranch: self.Parent.LocalName(),
		})
	}
	return nil
}
