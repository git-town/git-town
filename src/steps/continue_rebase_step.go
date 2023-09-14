package steps

// ContinueRebaseStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseStep struct {
	EmptyStep
}

func (step *ContinueRebaseStep) CreateAbortSteps() []Step {
	return []Step{&AbortRebaseStep{}}
}

func (step *ContinueRebaseStep) CreateContinueSteps() []Step {
	return []Step{step}
}

func (step *ContinueRebaseStep) Run(args RunArgs) error {
	hasRebaseInProgress, err := args.Run.Backend.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if hasRebaseInProgress {
		return args.Run.Frontend.ContinueRebase()
	}
	return nil
}
