package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// CommitRevert adds a commit to the current branch
// that reverts the commit with the given SHA.
type CommitRevert struct {
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitRevert) Run(args shared.RunArgs) error {
	return args.Git.RevertCommit(args.Frontend, self.SHA)
}
