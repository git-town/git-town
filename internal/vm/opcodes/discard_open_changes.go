package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

// DiscardOpenChanges resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChanges struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *DiscardOpenChanges) Run(args shared.RunArgs) error {
	return args.Git.DiscardOpenChanges(args.Frontend)
}
