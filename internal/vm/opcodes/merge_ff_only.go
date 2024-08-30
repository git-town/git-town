package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// SquashMerge squash merges the branch with the given name into the current branch.
type MergeFastForward struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *MergeFastForward) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
	}
}

func (self *MergeFastForward) CreateAutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *MergeFastForward) Run(args shared.RunArgs) error {
	return args.Git.MergeFastForward(args.Frontend, self.Branch)
}

func (self *MergeFastForward) ShouldAutomaticallyUndoOnError() bool {
	return true
}
