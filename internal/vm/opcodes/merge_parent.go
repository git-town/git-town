package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// MergeParent merges the given parent branch into the current branch.
type MergeParent struct {
	Parent                  gitdomain.BranchName // the currently active parent, after all remotely deleted parents were removed
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParent) Run(args shared.RunArgs) error {
	return args.Git.MergeBranchNoEdit(args.Frontend, self.Parent)
}
