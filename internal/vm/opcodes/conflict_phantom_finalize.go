package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type ConflictPhantomFinalize struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomFinalize) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *ConflictPhantomFinalize) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *ConflictPhantomFinalize) Run(args shared.RunArgs) error {
	unmergedFiles, err := args.Git.FileConflictQuickInfos(args.Backend)
	if err != nil {
		return err
	}
	if len(unmergedFiles) > 0 {
		// there are still unmerged files --> these are not phantom merge conflicts, let the user sort this out
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
