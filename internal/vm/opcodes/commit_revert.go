package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CommitRevert adds a commit to the current branch
// that reverts the commit with the given SHA.
type CommitRevert struct {
	SHA gitdomain.SHA
}

func (self *CommitRevert) Run(args shared.RunArgs) error {
	return args.Git.RevertCommit(args.Frontend, self.SHA)
}
