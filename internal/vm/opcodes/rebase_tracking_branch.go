package opcodes

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

// RebaseTrackingBranch rebases the current feature branch against its tracking branch.
type RebaseTrackingBranch struct {
	CurrentBranch           gitdomain.LocalBranchName
	PushBranches            configdomain.PushBranches
	RemoteBranch            gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseTrackingBranch) Run(args shared.RunArgs) error {
	inSync, err := args.Git.BranchInSyncWithTracking(args.Backend, self.CurrentBranch, args.Config.Value.NormalConfig.DevRemote)
	if err != nil {
		return err
	}
	if inSync {
		return nil
	}
	opcodes := []shared.Opcode{
		&RebaseBranch{
			Branch:                  self.RemoteBranch.BranchName(),
			undeclaredOpcodeMethods: undeclaredOpcodeMethods{},
		},
	}
	if self.PushBranches {
		opcodes = append(opcodes, &PushCurrentBranchForceIfNeeded{
			ForceIfIncludes: true,
		})
	}
	args.PrependOpcodes(opcodes...)
	return nil
}
