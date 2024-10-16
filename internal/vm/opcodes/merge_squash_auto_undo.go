package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type MergeSquashAutoUndo struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeSquashAutoUndo) AbortProgram() []shared.Opcode {
	result := []shared.Opcode{
		&ChangesDiscard{},
	}
	return result
}

func (self *MergeSquashAutoUndo) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *MergeSquashAutoUndo) Run(args shared.RunArgs) error {
	return args.Git.SquashMerge(args.Frontend, self.Branch)
}

func (self *MergeSquashAutoUndo) ShouldUndoOnError() bool {
	return true
}
