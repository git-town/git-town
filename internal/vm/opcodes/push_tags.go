package opcodes

import "github.com/git-town/git-town/v18/internal/vm/shared"

// PushTags pushes newly created Git tags to origin.
type PushTags struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushTags) Run(args shared.RunArgs) error {
	return args.Git.PushTags(args.Frontend, args.Config.Value.NormalConfig.NoPushHook())
}
