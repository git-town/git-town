package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// MergeContinue finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type MergeContinue struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeContinue) Run(args shared.RunArgs) error {
	if args.Git.HasMergeInProgress(args.Backend) {
		args.PrependOpcodes(&Commit{
			AuthorOverride:                 None[gitdomain.Author](),
			FallbackToDefaultCommitMessage: true,
			Message:                        None[gitdomain.CommitMessage](),
		})
	}
	return nil
}
