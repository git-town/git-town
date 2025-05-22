package opcodes

import (
	"fmt"
	"os"
	"time"

	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/subshell"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// rebases the current branch against the target branch while executing "git town swap", while moving the target branch onto the Onto branch.
type RebaseOntoKeepDeleted struct {
	BranchToRebaseOnto      gitdomain.BranchName
	CommitsToRemove         gitdomain.Location
	Upstream                Option[gitdomain.LocalBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOntoKeepDeleted) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoKeepDeleted) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoKeepDeleted) Run(args shared.RunArgs) error {
	// wait here in tests
	if len(os.Getenv(subshell.TestToken)) > 0 {
		time.Sleep(1 * time.Second)
	}
	err := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove, self.Upstream)
	if err != nil {
		conflictingFiles, err := args.Git.FileConflictQuickInfos(args.Backend)
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
		err = args.Git.ContinueRebase(args.Frontend)
		if err != nil {
			return fmt.Errorf("cannot continue rebase: %w", err)
		}
	}
	return nil
}
