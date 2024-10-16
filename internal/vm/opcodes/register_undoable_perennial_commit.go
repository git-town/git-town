package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
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
