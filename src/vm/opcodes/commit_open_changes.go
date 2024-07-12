package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/vm/shared"
)

// CommitOpenChanges commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChanges struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitOpenChanges) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{self}
}

func (self *CommitOpenChanges) Run(args shared.RunArgs) error {
	err := args.Git.StageFiles(args.Frontend, "-A")
	if err != nil {
		return err
	}
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	return args.Git.CommitStagedChanges(args.Frontend, fmt.Sprintf("WIP on %s", currentBranch))
}
