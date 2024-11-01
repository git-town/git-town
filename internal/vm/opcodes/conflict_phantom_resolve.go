package opcodes

import (
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConflictPhantomResolve struct {
	FilePath                string
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
	err := args.Git.CheckoutOurVersion(args.Frontend, self.FilePath)
	if err != nil {
		return err
	}
	err = args.Git.StageFiles(args.Frontend, self.FilePath)
	if err != nil {
		return err
	}
	return nil
}
