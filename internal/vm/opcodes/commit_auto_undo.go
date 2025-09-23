package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// CommitAutoUndo is a Commit that automatically aborts the Git Town command on failure.
type CommitAutoUndo struct {
	AuthorOverride                 Option[gitdomain.Author]
	FallbackToDefaultCommitMessage bool
	Message                        Option[gitdomain.CommitMessage]
}

func (self *CommitAutoUndo) Abort() []shared.Opcode {
	return []shared.Opcode{
		&ChangesDiscard{},
	}
}

func (self *CommitAutoUndo) AutomaticUndoError() error {
	return errors.New(messages.ShipExitMergeError)
}

func (self *CommitAutoUndo) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, configdomain.UseMessageWithFallbackToDefault(self.Message, self.FallbackToDefaultCommitMessage), self.AuthorOverride, configdomain.CommitHookEnabled)
}
