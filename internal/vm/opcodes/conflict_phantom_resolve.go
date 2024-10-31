package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

type ConflictPhantomResolve struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomResolve) Run(args shared.RunArgs) error {
	args.Git.StageFiles(args.Frontend)
	return nil
}
