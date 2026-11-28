package opcodes

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
)

// BranchPushesRemove removes queued force pushes for the given branch.
type BranchPushesRemove struct {
	Branch gitdomain.LocalBranchName
}

func (self *BranchPushesRemove) Run(args shared.RunArgs) error {
	args.RemoveOpcodes(func(opcode shared.Opcode) bool {
		push, isPush := opcode.(*PushCurrentBranchForceIfNeeded)
		return isPush && push.CurrentBranch == self.Branch
	})
	return nil
}
