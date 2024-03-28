package opcodes

import (
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// CommitOpenChanges commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitSquashedChanges struct {
	OldCommitMessages []string
	undeclaredOpcodeMethods
}

func (self *CommitSquashedChanges) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.Commit("", "")
}
