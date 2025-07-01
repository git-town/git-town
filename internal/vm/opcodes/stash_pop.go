package opcodes

import (
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// restores stashed away changes into the workspace
type StashPop struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *StashPop) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{&StashDrop{}}
}

func (self *StashPop) Run(args shared.RunArgs) error {
	if err := args.Git.PopStash(args.Frontend); err != nil {
		args.FinalMessages.Add(messages.DiffConflictWithMain)
	}
	return nil
}
