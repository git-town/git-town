package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// MergeFastForward fast-forwards the current branch to point to the commits on the given branch.
type MergeFastForward struct {
	Branch gitdomain.BranchName
}

func (self *MergeFastForward) Abort() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *MergeFastForward) AutomaticUndoError() error {
	return errors.New(messages.ShipExitMergeError)
}

func (self *MergeFastForward) Run(args shared.RunArgs) error {
	return args.Git.MergeFastForward(args.Frontend, self.Branch)
}

func (self *MergeFastForward) ShouldUndoOnError() bool {
	return true
}
