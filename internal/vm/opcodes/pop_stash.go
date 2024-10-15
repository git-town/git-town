package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type PopStash struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PopStash) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{&StashDrop{}}
}

func (self *PopStash) Run(args shared.RunArgs) error {
	err := args.Git.PopStash(args.Frontend)
	if err != nil {
		return errors.New(messages.DiffConflictWithMain)
	}
	return nil
}
