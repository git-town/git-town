package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ConflictMergePhantomFinalize struct{}

func (self *ConflictMergePhantomFinalize) Abort() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *ConflictMergePhantomFinalize) Continue() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *ConflictMergePhantomFinalize) Run(args shared.RunArgs) error {
	unresolvedFiles, err := args.Git.FileConflicts(args.Backend)
	if err != nil {
		return err
	}
	if len(unresolvedFiles) > 0 {
		// there are still unresolved files --> these are not phantom merge conflicts, let the user sort this out
		return errors.New(messages.ConflictMerge)
	}
	// here all merge conflicts have been resolved --> commit to finish the merge conflict and continue the program
	args.PrependOpcodes(
		&Commit{
			AuthorOverride:                 None[gitdomain.Author](),
			FallbackToDefaultCommitMessage: true,
			Message:                        None[gitdomain.CommitMessage](),
		},
	)
	return nil
}
