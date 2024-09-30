package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type RegisterUndoablePerennialCommit struct {
	Branch gitdomain.BranchName
	undeclaredOpcodeMethods
}

func (self *RegisterUndoablePerennialCommit) Run(args shared.RunArgs) error {
	squashedCommitSHA, err := args.Git.SHAForBranch(args.Backend, self.Branch)
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}
