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
	// We could check whether the branch is in sync with its tracking branch here
	// and then not sync with the tracking branch.
	// But that leads to flaky tests where tests sometimes sync with the tracking branch and sometimes not.
	// This is due to timestamps having only 1-second resolution.
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
