package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// RevertCommitIfNeeded adds a commit to the current branch
// that reverts the commit with the given SHA.
type RevertCommit struct {
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RevertCommit) Run(args shared.RunArgs) error {
	return args.Git.RevertCommit(args.Frontend, self.SHA)
}
