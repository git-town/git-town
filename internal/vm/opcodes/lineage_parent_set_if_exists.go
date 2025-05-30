package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// LineageParentSetIfExists sets the given parent branch as the parent of the given branch,
// but only the latter exists.
type LineageParentSetIfExists struct {
	Branch                  gitdomain.LocalBranchName
	Parent                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageParentSetIfExists) Run(args shared.RunArgs) error {
	if !args.Git.BranchExists(args.Backend, self.Branch) {
		return nil
	}
	args.PrependOpcodes(&LineageParentSet{Branch: self.Branch, Parent: self.Parent})
	return nil
}
