package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// registers the commit on the current perennial branch as undoable
type RegisterUndoablePerennialCommit struct {
	Parent                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RegisterUndoablePerennialCommit) Run(args shared.RunArgs) error {
	squashedCommitSHA, err := args.Git.SHAForBranch(args.Backend, self.Parent)
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}
