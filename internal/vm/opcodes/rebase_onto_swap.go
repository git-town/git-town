package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// rebases the current branch against the target branch, while moving the target branch onto the Onto branch.
type RebaseOntoSwap struct {
	BranchToRebaseAgainst   gitdomain.BranchName
	BranchToRebaseOnto      gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOntoSwap) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoSwap) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoSwap) Run(args shared.RunArgs) error {
	err := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseAgainst, self.BranchToRebaseOnto, None[gitdomain.LocalBranchName]())
	if err != nil {
		conflictingFiles, err := args.Git.FileConflictQuickInfos(args.Backend)
		if err != nil {
			return fmt.Errorf("cannot determine conflicting files after rebase: %w", err)
		}
		for _, conflictingFile := range conflictingFiles {
			if conflictingChange, has := conflictingFile.CurrentBranchChange.Get(); has {
				_ = args.Git.CheckoutTheirVersion(args.Frontend, conflictingChange.FilePath)
				_ = args.Git.StageFiles(args.Frontend, conflictingChange.FilePath)
			} else if baseChange, has := conflictingFile.BaseChange.Get(); has {
				_ = args.Git.StageFiles(args.Frontend, baseChange.FilePath)
			}
		}
		_ = args.Git.ContinueRebase(args.Frontend)
	}
	return nil
}
