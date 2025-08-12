package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// rebases the current branch against its parent, which exists only remotely
type RebaseAncestorRemote struct {
	Branch                  gitdomain.LocalBranchName
	Parent                  gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseAncestorRemote) Run(args shared.RunArgs) error {
	isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, self.Parent.BranchName())
	if err != nil {
		return err
	}
	if !isInSync {
		args.PrependOpcodes(
			&RebaseBranch{
				Branch: self.Parent.BranchName(),
			},
		)
	}
	return nil
}
