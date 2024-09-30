package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ResetRemoteBranchToSHAIfNeeded sets the given remote branch to the given SHA,
// but only if it currently has a particular SHA.
type ResetRemoteBranchToSHAIfNeeded struct {
	Branch                  gitdomain.RemoteBranchName
	MustHaveSHA             gitdomain.SHA
	SetToSHA                gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ResetRemoteBranchToSHAIfNeeded) Run(args shared.RunArgs) error {
	currentSHA, err := args.Git.SHAForBranch(args.Backend, self.Branch.BranchName())
	if err != nil {
		return err
	}
	if currentSHA != self.MustHaveSHA {
		return fmt.Errorf(messages.BranchHasWrongSHA, self.Branch, self.SetToSHA, self.MustHaveSHA, currentSHA)
	}
	args.PrependOpcodes(&ResetRemoteBranchToSHA{Branch: self.Branch, SetToSHA: self.SetToSHA})
	return nil
}
