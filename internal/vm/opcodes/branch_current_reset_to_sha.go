package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchCurrentResetToSHA undoes all commits on the current branch
// all the way until the given SHA.
type BranchCurrentResetToSHA struct {
	Hard                    bool
	SetToSHA                gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchCurrentResetToSHA) Run(args shared.RunArgs) error {
	return args.Git.ResetCurrentBranchToSHA(args.Frontend, self.SetToSHA, self.Hard)
}
