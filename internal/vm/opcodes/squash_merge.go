package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type SquashMerge struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SquashMerge) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&DiscardOpenChanges{},
	}
}

func (self *SquashMerge) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *SquashMerge) Run(args shared.RunArgs) error {
	return args.Git.SquashMerge(args.Frontend, self.Branch)
}

func (self *SquashMerge) ShouldUndoOnError() bool {
	return true
}
