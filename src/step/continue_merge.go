package step

// ContinueMerge finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMerge struct {
	Empty
}

func (step *ContinueMerge) CreateContinueProgram() []Step {
	return []Step{
		step,
	}
}

func (step *ContinueMerge) Run(args RunArgs) error {
	if args.Runner.Backend.HasMergeInProgress() {
		return args.Runner.Frontend.CommitNoEdit()
	}
	return nil
}
