package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// MergeFastForward fast-forwards the current branch to point to the commits on the given branch.
type MergeFastForward struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *MergeFastForward) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *MergeFastForward) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *MergeFastForward) Run(args shared.RunArgs) error {
	return args.Git.MergeFastForward(args.Frontend, self.Branch)
}

func (self *MergeFastForward) ShouldUndoOnError() bool {
	return true
}
