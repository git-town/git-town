package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// RebaseOntoRemoveDeleted rebases the current branch against the target branch, while moving the target branch onto the Onto branch.
// If there are merge conflicts,
type RebaseOntoRemoveDeleted struct {
	BranchToRebaseOnto gitdomain.LocalBranchName
	CommitsToRemove    gitdomain.BranchName
}

func (self *RebaseOntoRemoveDeleted) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoRemoveDeleted) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoRemoveDeleted) Run(args shared.RunArgs) error {
	err := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove.Location())
	if err != nil || args.Config.Value.NormalConfig.AutoResolve.NoAutoResolve() {
		return err
	}
	// Here the rebase-onto has failed.
	// The branch that gets rebased onto will be deleted.
	// We therefore don't need to bother the user with resolving the merge conflict
	// and can resolve it ourselves.
	conflictingFiles, err := args.Git.FileConflicts(args.Backend)
	if err != nil {
		return fmt.Errorf("cannot determine conflicting files after rebase: %w", err)
	}
	opcodes := []shared.Opcode{}
	for _, conflictingFile := range conflictingFiles {
		if conflictingChange, has := conflictingFile.CurrentBranchChange.Get(); has {
			opcodes = append(opcodes,
				&ConflictResolve{
					FilePath:   conflictingChange.FilePath,
					Resolution: gitdomain.ConflictResolutionTheirs,
				},
				&FileStage{
					FilePath: conflictingChange.FilePath,
				},
			)
		} else if baseChange, has := conflictingFile.BaseChange.Get(); has {
			opcodes = append(opcodes, &FileRemove{
				FilePath: baseChange.FilePath,
			})
		}
	}
	opcodes = append(opcodes, &RebaseContinue{})
	args.PrependOpcodes(opcodes...)
	return nil
}
