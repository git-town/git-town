package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchCurrentResetToSHA undoes all commits on the current branch
// all the way until the given SHA.
type BranchCurrentResetToSHA struct {
	SHA gitdomain.SHA
}

func (self *BranchCurrentResetToSHA) Run(args shared.RunArgs) error {
	return args.Git.ResetCurrentBranchToSHA(args.Frontend, self.SHA)
}
