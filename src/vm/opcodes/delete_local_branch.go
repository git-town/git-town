package opcodes

import (
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// DeleteLocalBranch deletes the branch with the given name.
type DeleteLocalBranch struct {
	Branch gitdomain.LocalBranchName
	Force  bool
	undeclaredOpcodeMethods
}

func (self *DeleteLocalBranch) Run(args shared.RunArgs) error {
	useForce := self.Force
	if !useForce {
		parent := args.Lineage.Parent(self.Branch)
		hasUnmergedCommits, err := args.Runner.Backend.BranchHasUnmergedCommits(self.Branch, parent.Location())
		if err != nil {
			return err
		}
		if hasUnmergedCommits {
			useForce = true
		}
	}
	return args.Runner.Frontend.DeleteLocalBranch(self.Branch, useForce)
}
