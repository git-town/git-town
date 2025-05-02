package opcodes

import "github.com/git-town/git-town/v20/internal/vm/shared"

type ChangesUnstageAll struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ChangesUnstageAll) Run(args shared.RunArgs) error {
	return args.Git.UnstageAll(args.Frontend)
}
