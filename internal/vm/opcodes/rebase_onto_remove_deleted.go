package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// RebaseOntoRemoveDeleted rebases the current branch against the target branch, while moving the target branch onto the Onto branch.
// If there are merge conflicts,
type RebaseOntoRemoveDeleted struct {
	BranchToRebaseOnto      gitdomain.LocalBranchName
	CommitsToRemove         gitdomain.BranchName
	Upstream                Option[gitdomain.LocalBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOntoRemoveDeleted) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoRemoveDeleted) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoRemoveDeleted) Run(args shared.RunArgs) error {
	if err := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove.Location(), self.Upstream); err != nil {
		// Here the rebase-onto has failed.
		// The branch that gets rebased onto will be deleted.
		// We therefore don't need to bother the user with resolving the merge conflict
		// and can resolve it ourselves.
		conflictingFiles, err := args.Git.FileConflictQuickInfos(args.Backend)
		if err != nil {
			return fmt.Errorf("cannot determine conflicting files after rebase: %w", err)
		}
		for _, conflictingFile := range conflictingFiles {
			if conflictingChange, has := conflictingFile.CurrentBranchChange.Get(); has {
				_ = args.Git.ResolveConflict(args.Frontend, conflictingChange.FilePath, gitdomain.ConflictResolutionTheirs)
				_ = args.Git.StageFiles(args.Frontend, conflictingChange.FilePath)
			} else if baseChange, has := conflictingFile.BaseChange.Get(); has {
				_ = args.Git.RemoveFile(args.Frontend, baseChange.FilePath)
			}
		}
		_ = args.Git.ContinueRebase(args.Frontend)
	}
	return nil
}
