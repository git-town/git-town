package opcodes

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// BranchRemoteSetToSHA sets the given remote branch to the given SHA.
type BranchLocalSetToSHA struct {
	SetToSHA                gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalSetToSHA) Run(args shared.RunArgs) error {
	return args.Git.ResetCurrentBranchToSHA(args.Frontend, self.SetToSHA, true)
}
