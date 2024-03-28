package opcodes

import (
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// CommitOpenChanges commits all open changes as a new commit.
type CommitSquashedChanges struct {
	undeclaredOpcodeMethods
}

func (self *CommitSquashedChanges) Run(args shared.RunArgs) error {
	err := args.Runner.Frontend.StageFiles("-A")
	if err != nil {
		return err
	}
	return args.Runner.Frontend.Commit("", "")
}
