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
	if args.Run.Backend.HasMergeInProgress() {
		return args.Run.Frontend.CommitNoEdit()
	}
	return nil
}
