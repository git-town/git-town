package opcodes

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// rebases the current branch against the target branch while executing "git town swap", while moving the target branch onto the Onto branch.
type RebaseOntoKeepDeleted struct {
	BranchToRebaseOnto      gitdomain.BranchName
	CommitsToRemove         gitdomain.Location
	Upstream                Option[gitdomain.LocalBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOntoKeepDeleted) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoKeepDeleted) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoKeepDeleted) Run(args shared.RunArgs) error {
	// Fix for https://github.com/git-town/git-town/issues/4942.
	// Waiting here in end-to-end tests to ensure new timestamps for the rebased commits,
	// which avoids flaky end-to-end tests.
	if subshell.IsInTest() {
		time.Sleep(1 * time.Second)
	}
	if err := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove, self.Upstream); err != nil {
		conflictingFiles, err := args.Git.FileConflicts(args.Backend)
		if err != nil {
			return fmt.Errorf("cannot determine conflicting files after rebase: %w", err)
		}
		for _, conflictingFile := range conflictingFiles {
			if conflictingChange, has := conflictingFile.CurrentBranchChange.Get(); has {
				_ = args.Git.ResolveConflict(args.Frontend, conflictingChange.FilePath, gitdomain.ConflictResolutionTheirs)
				_ = args.Git.StageFiles(args.Frontend, conflictingChange.FilePath)
			} else if baseChange, has := conflictingFile.BaseChange.Get(); has {
				_ = args.Git.StageFiles(args.Frontend, baseChange.FilePath)
			}
		}
		if err = args.Git.ContinueRebase(args.Frontend); err != nil {
			return fmt.Errorf("cannot continue rebase: %w", err)
		}
	}
	return nil
}
