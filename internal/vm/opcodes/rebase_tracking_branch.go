package opcodes

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
)

// RebaseTrackingBranch rebases the current feature branch against its tracking branch.
type RebaseTrackingBranch struct {
	Branch       gitdomain.LocalBranchName
	RemoteBranch gitdomain.RemoteBranchName
}

func (self *RebaseTrackingBranch) Run(args shared.RunArgs) error {
	if args.Git.RemoteTipInReflog(args.Backend, self.Branch, self.RemoteBranch) {
		// The remote branch only has commits that the local branch used to have, i.e. the local branch was rewritten.
		// Rebasing onto it would replay the rewritten commits on top of their old versions.
		return nil
	}
	args.PrependOpcodes(
		&RebaseBranch{Branch: self.RemoteBranch.BranchName()},
	)
	return nil
}
