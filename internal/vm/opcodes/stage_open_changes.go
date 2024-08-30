package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

// CommitOpenChanges commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type StageOpenChanges struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *StageOpenChanges) Run(args shared.RunArgs) error {
	return args.Git.StageFiles(args.Frontend, "-A")
}
