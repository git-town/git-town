package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// deletes the given branch including all commits
type BranchLocalDeleteContent struct {
	BranchToRebaseOnto      gitdomain.BranchName
	BranchToDelete          gitdomain.LocalBranchName
	BranchToBeOn            gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalDeleteContent) Run(args shared.RunArgs) error {
	args.PrependOpcodes(
		&CheckoutIfNeeded{Branch: self.BranchToBeOn},
		&RebaseOnto{BranchToRebaseOnto: self.BranchToRebaseOnto, BranchToRebaseAgainst: self.BranchToDelete},
		&BranchLocalDelete{Branch: self.BranchToDelete},
	)
	return nil
}
