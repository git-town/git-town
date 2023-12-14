package opcode

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// PushCurrentBranch pushes the current branch to its existing tracking branch.
type PushCurrentBranch struct {
	CurrentBranch domain.LocalBranchName
	NoPushHook    configdomain.NoPushHook
	undeclaredOpcodeMethods
}

func (self *PushCurrentBranch) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *PushCurrentBranch) Run(args shared.RunArgs) error {
	shouldPush, err := args.Runner.Backend.ShouldPushBranch(self.CurrentBranch, self.CurrentBranch.TrackingBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Runner.GitTown.DryRun {
		return nil
	}
	return args.Runner.Frontend.PushCurrentBranch(self.NoPushHook)
}
