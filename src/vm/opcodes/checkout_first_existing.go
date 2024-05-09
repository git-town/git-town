package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// CheckoutIfExists does the same as Checkout
// but only if that branch actually exists.
type CheckoutFirstExisting struct {
	Branches                gitdomain.LocalBranchNames
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CheckoutFirstExisting) Run(args shared.RunArgs) error {
	candidates := args.BranchInfos.WithNames(self.Branches...).WithStatuses(gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusDeletedAtRemote).Names()
	candidates = append(candidates, args.Config.Config.MainBranch)
	branch := candidates[0]
	return (&Checkout{Branch: branch}).Run(args)
}
