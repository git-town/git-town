package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// ResetRemoteBranchToSHA sets the given remote branch to the given SHA,
// but only if it currently has a particular SHA.
type ResetRemoteBranchToSHA struct {
	Branch                  gitdomain.RemoteBranchName
	MustHaveSHA             gitdomain.SHA
	SetToSHA                gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ResetRemoteBranchToSHA) Run(args shared.RunArgs) error {
	currentSHA, err := args.Git.SHAForBranch(args.Backend, self.Branch.BranchName())
	if err != nil {
		return err
	}
	if currentSHA != self.MustHaveSHA {
		return fmt.Errorf(messages.BranchHasWrongSHA, self.Branch, self.SetToSHA, self.MustHaveSHA, currentSHA)
	}
	return args.Git.ResetRemoteBranchToSHA(args.Frontend, self.Branch, self.SetToSHA)
}
