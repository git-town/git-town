package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// BranchRemoteSetToSHA sets the given remote branch to the given SHA.
type BranchRemoteSetToSHA struct {
	Branch   gitdomain.RemoteBranchName
	SetToSHA gitdomain.SHA
}

func (self *BranchRemoteSetToSHA) Run(args shared.RunArgs) error {
	return args.Git.ResetRemoteBranchToSHA(args.Frontend, self.Branch, self.SetToSHA)
}
