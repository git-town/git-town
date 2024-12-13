package opcodes

import (
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// PushCurrentBranchForceIfNeeded force-pushes the branch with the given name to the origin remote.
type PushCurrentBranchForceIfNeeded struct {
	ForceIfIncludes         bool
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushCurrentBranchForceIfNeeded) Run(args shared.RunArgs) error {
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	shouldPush, err := args.Git.ShouldPushBranch(args.Backend, currentBranch)
	if err != nil {
		return err
	}
	if !shouldPush {
		return nil
	}
	args.PrependOpcodes(&PushCurrentBranchForce{ForceIfIncludes: self.ForceIfIncludes})
	return nil
}
