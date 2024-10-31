package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

type ConflictPhantomResolve struct {
	FilePath                string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomResolve) Run(args shared.RunArgs) error {
	// TODO
	args.Git.StageFiles(args.Frontend)
	return nil
}
