package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// CommitAutoUndo is a Commit that automatically aborts the Git Town command on failure.
type CommitAutoUndo struct {
	AuthorOverride                 Option[gitdomain.Author]
	FallbackToDefaultCommitMessage bool
	Message                        Option[gitdomain.CommitMessage]
	undeclaredOpcodeMethods        `exhaustruct:"optional"`
}

func (self *CommitAutoUndo) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&ChangesDiscard{},
	}
}

func (self *CommitAutoUndo) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *CommitAutoUndo) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, configdomain.UseMessageWithFallbackToDefault(self.Message, self.FallbackToDefaultCommitMessage), self.AuthorOverride, configdomain.CommitHookEnabled)
}

func (self *CommitAutoUndo) ShouldUndoOnError() bool {
	return true
}
