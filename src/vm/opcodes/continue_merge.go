package opcodes

import "github.com/git-town/git-town/v14/src/vm/shared"

// ContinueMerge finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMerge struct {
	undeclaredOpcodeMethods
}

func (self *ContinueMerge) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *ContinueMerge) Run(args shared.RunArgs) error {
	if args.Runner.Backend.HasMergeInProgress() {
		return args.Runner.Frontend.CommitNoEdit()
	}
	return nil
}
