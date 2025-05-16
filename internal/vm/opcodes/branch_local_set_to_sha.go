package opcodes

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// BranchLocalSetToSHA sets the given local branch to the given SHA.
type BranchLocalSetToSHA struct {
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalSetToSHA) Run(args shared.RunArgs) error {
	return args.Git.ResetCurrentBranchToSHA(args.Frontend, self.SHA)
}
