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
	repoStatus, err := args.Runner.Backend.RepoStatus()
	if err != nil {
		return err
	}
	if repoStatus.RebaseInProgress {
		return args.Runner.Frontend.ContinueRebase()
	}
	return nil
}
