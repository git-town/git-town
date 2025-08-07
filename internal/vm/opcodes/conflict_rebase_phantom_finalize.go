package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ConflictRebasePhantomFinalize struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictRebasePhantomFinalize) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *ConflictRebasePhantomFinalize) Run(args shared.RunArgs) error {
	unresolvedFiles, err := args.Git.FileConflicts(args.Backend)
	if err != nil {
		return err
	}
	if len(unresolvedFiles) > 0 {
		// there are still unresolved files --> these are not phantom conflicts, let the user sort this out
		return errors.New(messages.ConflictRebase)
	}
	// here all rebase conflicts have been resolved --> finish the rebase and continue the program
	args.PrependOpcodes(
		&RebaseContinueIfNeeded{},
	)
	return nil
}
