package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// PushCurrentBranch pushes the current branch to its existing tracking branch.
type PushCurrentBranch struct {
	CurrentBranch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *PushCurrentBranch) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *PushCurrentBranch) Run(args shared.RunArgs) error {
	shouldPush, err := args.Backend.ShouldPushBranch(self.CurrentBranch, self.CurrentBranch.TrackingBranch())
	if err != nil {
		return err
	}
	if !shouldPush {
		return nil
	}
	return args.Frontend.PushCurrentBranch(args.Config.Config.NoPushHook())
}
