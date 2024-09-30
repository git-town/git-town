package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// PushCurrentBranchIfNeeded pushes the current branch to its existing tracking branch
// if it has unpushed commits.
type PushCurrentBranchIfNeeded struct {
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushCurrentBranchIfNeeded) Run(args shared.RunArgs) error {
	shouldPush, err := args.Git.ShouldPushBranch(args.Backend, self.CurrentBranch)
	if err != nil {
		return err
	}
	if !shouldPush {
		return nil
	}
	args.PrependOpcodes(&PushCurrentBranch{CurrentBranch: self.CurrentBranch})
	return nil
}
