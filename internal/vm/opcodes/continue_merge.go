package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// ContinueMerge finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMerge struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ContinueMerge) Run(args shared.RunArgs) error {
	if args.Git.HasMergeInProgress(args.Backend) {
		return args.Git.CommitNoUserInput(args.Frontend, None[gitdomain.CommitMessage]())
	}
	return nil
}
