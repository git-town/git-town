package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// ResetRemoteBranchToSHA sets the given remote branch to the given SHA,
// but only if it currently has a particular SHA.
type ResetRemoteBranchToSHA struct {
	Branch      gitdomain.RemoteBranchName
	MustHaveSHA gitdomain.SHA
	SetToSHA    gitdomain.SHA
	undeclaredOpcodeMethods
}

func (self *ResetRemoteBranchToSHA) Run(args shared.RunArgs) error {
	currentSHA, err := args.Runner.Backend.SHAForBranch(self.Branch.BranchName())
	if err != nil {
		return err
	}
	if currentSHA != self.MustHaveSHA {
		return fmt.Errorf(messages.BranchHasWrongSHA, self.Branch, self.SetToSHA, self.MustHaveSHA, currentSHA)
	}
	return args.Runner.Frontend.ResetRemoteBranchToSHA(self.Branch, self.SetToSHA)
}
