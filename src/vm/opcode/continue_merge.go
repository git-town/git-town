package opcode

import "github.com/git-town/git-town/v9/src/vm/shared"

// ContinueMerge finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMerge struct {
	undeclaredOpcodeMethods
}

func (step *ContinueMerge) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		step,
	}
}

func (step *ContinueMerge) Run(args shared.RunArgs) error {
	if args.Runner.Backend.HasMergeInProgress() {
		return args.Runner.Frontend.CommitNoEdit()
	}
	return nil
}
