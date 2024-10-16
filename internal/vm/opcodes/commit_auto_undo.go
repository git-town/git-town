package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// CommitAutoUndo is a Commit that automatically aborts the Git Town command on failure.
type CommitAutoUndo struct {
	AuthorOverride                 Option[gitdomain.Author]
	FallbackToDefaultCommitMessage bool
	Message                        Option[gitdomain.CommitMessage]
	undeclaredOpcodeMethods        `exhaustruct:"optional"`
}

func (self *CommitAutoUndo) AbortProgram() []shared.Opcode {
	result := []shared.Opcode{
		&ChangesDiscard{},
	}
	return result
}

func (self *CommitAutoUndo) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *CommitAutoUndo) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, self.Message, self.FallbackToDefaultCommitMessage, self.AuthorOverride)
}

func (self *CommitAutoUndo) ShouldUndoOnError() bool {
	return true
}
