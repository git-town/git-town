package opcodes

import (
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

// PushCurrentBranchIfNeeded pushes the current branch to its existing tracking branch
// if it has unpushed commits.
type PushCurrentBranchIfNeeded struct {
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushCurrentBranchIfNeeded) Run(args shared.RunArgs) error {
	// check if branch still exists
	// the branch could not exist at this point if it was pruned at runtime due to being empty
	branchExists := args.Git.BranchExists(args.Backend, self.CurrentBranch)
	if !branchExists {
		return nil
	}
	shouldPush, err := args.Git.ShouldPushBranch(args.Backend, self.CurrentBranch, args.Config.Value.NormalConfig.DevRemote)
	if err != nil {
		return err
	}
	if !shouldPush {
		return nil
	}
	args.PrependOpcodes(&PushCurrentBranch{})
	return nil
}
