package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

type StashOpenChanges struct{}

func (self *StashOpenChanges) Run(args shared.RunArgs) error {
	return args.Git.Stash(args.Frontend)
}
