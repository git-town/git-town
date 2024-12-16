package opcodes

import "github.com/git-town/git-town/v17/internal/vm/shared"

// PushTags pushes newly created Git tags to origin.
type PushTags struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushTags) Run(args shared.RunArgs) error {
	return args.Git.PushTags(args.Frontend)
}
