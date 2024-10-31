package opcodes

import (
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConflictPhantomResolve struct {
	FilePath                string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomResolve) Run(args shared.RunArgs) error {
	err := args.Git.CheckoutOurs(args.Frontend, self.FilePath)
	if err != nil {
		return err
	}
	err = args.Git.StageFiles(args.Frontend, self.FilePath)
	if err != nil {
		return err
	}
	return nil
}
