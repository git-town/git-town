package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchRemoteSetToSHA sets the given remote branch to the given SHA.
type BranchRemoteSetToSHA struct {
	Branch                  gitdomain.RemoteBranchName
	SetToSHA                gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchRemoteSetToSHA) Run(args shared.RunArgs) error {
	return args.Git.ResetRemoteBranchToSHA(args.Frontend, self.Branch, self.SetToSHA)
}
