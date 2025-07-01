package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ConflictPhantomResolve struct {
	FilePath                string
	Resolution              gitdomain.ConflictResolution
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomResolve) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *ConflictPhantomResolve) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *ConflictPhantomResolve) Run(args shared.RunArgs) error {
	if err := args.Git.ResolveConflict(args.Frontend, self.FilePath, self.Resolution); err != nil {
		return err
	}
	return args.Git.StageFiles(args.Frontend, self.FilePath)
}
