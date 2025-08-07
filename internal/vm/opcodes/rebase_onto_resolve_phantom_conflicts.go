package opcodes

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// rebases the current branch against the target branch while executing "git town swap", while moving the target branch onto the Onto branch.
type RebaseOntoResolvePhantomConflicts struct {
	BranchToRebaseOnto      gitdomain.BranchName
	CommitsToRemove         gitdomain.Location
	CurrentBranch           gitdomain.LocalBranchName
	Upstream                Option[gitdomain.LocalBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOntoResolvePhantomConflicts) Run(args shared.RunArgs) error {
	if subshell.IsInTest() {
		// Fix for https://github.com/git-town/git-town/issues/4942.
		// Waiting here in end-to-end tests to ensure new timestamps for the rebased commits,
		// which avoids flaky end-to-end tests.
		time.Sleep(1 * time.Second)
	}
	fmt.Println("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD BEFORE")
	if err := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove, self.Upstream); err != nil {
		fmt.Println("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD", err)
		args.PrependOpcodes(&ConflictRebasePhantomResolveAll{
			BranchToRebaseOnto: self.BranchToRebaseOnto,
			CurrentBranch:      self.CurrentBranch,
		})
	}
	return nil
}
