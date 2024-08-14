package opcodes

import (
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// CommitOpenChanges commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChanges struct {
	AddAll                  bool
	Message                 string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitOpenChanges) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{self}
}

func (self *CommitOpenChanges) Run(args shared.RunArgs) error {
	if self.AddAll {
		err := args.Git.StageFiles(args.Frontend, "-A")
		if err != nil {
			return err
		}
	}
	return args.Git.CommitStagedChanges(args.Frontend, self.Message)
}
