package steps

// ContinueMergeStep finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMergeStep struct {
	EmptyStep
}

func (step *ContinueMergeStep) CreateContinueSteps() []Step {
	return []Step{step}
}

func (step *ContinueMergeStep) Run(args RunArgs) error {
	if args.Runner.Backend.HasMergeInProgress() {
		return args.Runner.Frontend.CommitNoEdit()
	}
	return nil
}
