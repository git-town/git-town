package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// rebases a branch against its ancestor branch, which exists only remotely
type RebaseAncestorRemote struct {
	Ancestor                gitdomain.RemoteBranchName
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseAncestorRemote) Run(args shared.RunArgs) error {
	isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, self.Ancestor.BranchName())
	if err != nil {
		return err
	}
	if !isInSync {
		args.PrependOpcodes(
			&RebaseBranch{
				Branch: self.Ancestor.BranchName(),
			},
		)
	}
	return nil
}
